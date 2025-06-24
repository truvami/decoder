package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/truvami/decoder/internal/logger"
	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	nomadxlDecoder "github.com/truvami/decoder/pkg/decoder/nomadxl/v1"
	nomadxsDecoder "github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
	smartlabelDecoder "github.com/truvami/decoder/pkg/decoder/smartlabel/v1"
	tagslDecoder "github.com/truvami/decoder/pkg/decoder/tagsl/v1"
	tagxlDecoder "github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/encoder"
	nomadxsEncoder "github.com/truvami/decoder/pkg/encoder/nomadxs/v1"
	smartlabelEncoder "github.com/truvami/decoder/pkg/encoder/smartlabel/v1"
	tagslEncoder "github.com/truvami/decoder/pkg/encoder/tagsl/v1"
	"github.com/truvami/decoder/pkg/solver"
	"github.com/truvami/decoder/pkg/solver/aws"
	"github.com/truvami/decoder/pkg/solver/loracloud"
	"go.uber.org/zap"
)

var host string
var port uint16
var health bool
var metrics bool

func init() {
	httpCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the HTTP server to")
	httpCmd.Flags().Uint16Var(&port, "port", 8080, "Port to bind the HTTP server to")
	httpCmd.Flags().BoolVar(&health, "health", false, "Enable /health endpoint")
	httpCmd.Flags().BoolVar(&metrics, "metrics", false, "Enable prometheus /metrics endpoint")
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start the HTTP server for the decoder.",
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure logger is initialized
		if logger.Logger == nil {
			logger.NewLogger()
			defer logger.Sync()
		}

		router := http.NewServeMux()

		// health endpoint
		if health {
			router.HandleFunc("/health", healthHandler)
		}

		// metrics endpoint
		if metrics {
			logger.Logger.Debug("enabling prometheus metrics endpoint")
			router.Handle("/metrics", promhttp.Handler())
		}

		type decoderEndpoint struct {
			path    string
			decoder decoder.Decoder
		}

		ctx := context.Background()
		if cmd != nil {
			ctx = cmd.Context()
		}

		var solver solver.SolverV1
		var err error

		switch strings.ToLower(Solver) {
		case "aws":
			solver, err = aws.NewAwsPositionEstimateClient(ctx, logger.Logger)
			if err != nil {
				logger.Logger.Error("error while creating AWS position estimate client", zap.Error(err))
				os.Exit(1)
			}
		case "loracloud":
			if LoracloudAccessToken == "" {
				logger.Logger.Error("loracloud access token is required for loracloud solver")
				os.Exit(1)
			}
			solver = loracloud.NewLoracloudMiddleware(ctx, LoracloudAccessToken, logger.Logger)
		}

		var decoders []decoderEndpoint = []decoderEndpoint{
			{"tagsl/v1", tagslDecoder.NewTagSLv1Decoder(tagslDecoder.WithSkipValidation(SkipValidation))},
			{"tagxl/v1", tagxlDecoder.NewTagXLv1Decoder(ctx, solver, logger.Logger, tagxlDecoder.WithSkipValidation(SkipValidation))},
			{"nomadxs/v1", nomadxsDecoder.NewNomadXSv1Decoder(nomadxsDecoder.WithSkipValidation(SkipValidation))},
			{"nomadxl/v1", nomadxlDecoder.NewNomadXLv1Decoder(nomadxlDecoder.WithSkipValidation(SkipValidation))},
			{"smartlabel/v1", smartlabelDecoder.NewSmartLabelv1Decoder(ctx, solver, logger.Logger, smartlabelDecoder.WithSkipValidation(SkipValidation))},
		}

		// add the decoders
		for _, d := range decoders {
			addDecoder(ctx, router, d.path, d.decoder)
		}

		// Define encoder endpoints
		type encoderEndpoint struct {
			path    string
			encoder encoder.Encoder
		}

		var encoders []encoderEndpoint = []encoderEndpoint{
			{"encode/tagsl/v1", tagslEncoder.NewTagSLv1Encoder()},
			{"encode/nomadxs/v1", nomadxsEncoder.NewNomadXSv1Encoder()},
			{"encode/smartlabel/v1", smartlabelEncoder.NewSmartlabelv1Encoder()},
		}

		// add the encoders
		for _, e := range encoders {
			addEncoder(router, e.path, e.encoder)
		}

		// middleware
		handler := loggingMiddleware(logger.Logger, router)

		logger.Logger.Info("starting HTTP server", zap.String("host", host), zap.Uint64("port", uint64(port)))
		err = http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), handler)

		if err != nil {
			logger.Logger.Error("error while starting HTTP server", zap.Error(err))
			os.Exit(1)
		}
	},
}

func addDecoder(ctx context.Context, router *http.ServeMux, path string, decoder decoder.Decoder) {
	logger.Logger.Debug("adding decoder", zap.String("path", path))
	router.HandleFunc("POST /"+path, getHandler(ctx, decoder))
}

func getHandler(ctx context.Context, targetDecoder decoder.Decoder) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Port    uint8  `json:"port" validate:"required,gt=0,lte=255"`
			Payload string `json:"payload" validate:"required,hexadecimal"`
			DevEUI  string `json:"devEui" validate:"omitempty,hexadecimal,len=16"`
		}

		// decode the request
		var req request

		logger.Logger.Debug("decoding request")
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Logger.Error("error while decoding request", zap.Error(err))

			setBody(w, http.StatusBadRequest, map[string]any{
				"error": err.Error(),
				"docs":  "https://docs.truvami.com",
			})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			logger.Logger.Error("request validation failed", zap.Error(err))
			setBody(w, http.StatusBadRequest, map[string]any{
				"error": "request validation failed",
				"docs":  "https://docs.truvami.com",
			})
			return
		}

		logger.Logger.Debug("set context values",
			zap.String("devEui", req.DevEUI),
			zap.Uint8("port", req.Port),
			zap.String("payload", req.Payload),
		)
		ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, req.DevEUI)
		ctx = context.WithValue(ctx, decoder.PORT_CONTEXT_KEY, req.Port)
		ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 1) // Default frame count, can be adjusted as needed

		logger.Logger.Debug("decoding payload")

		var warnings []string = nil
		data, err := targetDecoder.Decode(ctx, req.Payload, req.Port)
		if err != nil {
			if errors.Is(err, helpers.ErrValidationFailed) {
				warnings = []string{}
				for _, err := range helpers.UnwrapError(err) {
					logger.Logger.Warn("validation error", zap.Error(err), zap.String("devEui", req.DevEUI), zap.Uint8("port", req.Port))
					warnings = append(warnings, err.Error())
				}
				logger.Logger.Warn("validation for some fields failed - are you using the correct port?")
			} else {
				logger.Logger.Error("error while decoding payload", zap.Error(err), zap.String("devEui", req.DevEUI), zap.Uint8("port", req.Port))

				setBody(w, http.StatusBadRequest, map[string]any{
					"error": err.Error(),
					"docs":  "https://docs.truvami.com",
				})
				return
			}
		}

		logger.Logger.Info("payload decoded successfully", zap.String("devEui", req.DevEUI), zap.Uint8("port", req.Port))
		setBody(w, http.StatusOK, map[string]any{
			"data":     data.Data,
			"warnings": warnings,
		})
	}
}

func addEncoder(router *http.ServeMux, path string, encoder encoder.Encoder) {
	logger.Logger.Debug("adding encoder", zap.String("path", path))
	router.HandleFunc("POST /"+path, getEncoderHandler(encoder))
}

func getEncoderHandler(encoder encoder.Encoder) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// First, decode the request to get the port and raw payload
		var rawReq struct {
			Port    uint8           `json:"port" validate:"required,gt=0,lte=255"`
			Payload json.RawMessage `json:"payload" validate:"required"`
			DevEUI  string          `json:"devEui" validate:"omitempty,hexadecimal,len=16"`
		}

		logger.Logger.Debug("decoding request")
		err := json.NewDecoder(r.Body).Decode(&rawReq)
		if err != nil {
			logger.Logger.Error("error while decoding request", zap.Error(err))

			setBody(w, http.StatusBadRequest, map[string]any{
				"error": err.Error(),
				"docs":  "https://docs.truvami.com",
			})
			return
		}

		if err := validator.New().Struct(rawReq); err != nil {
			logger.Logger.Error("request validation failed", zap.Error(err))
			setBody(w, http.StatusBadRequest, map[string]any{
				"error": "request validation failed",
				"docs":  "https://docs.truvami.com",
			})
			return
		}

		// Now we have the raw payload and port, we can determine the correct struct type
		// and unmarshal the payload into it
		var structPayload any

		// This is a simplified example - in a real implementation, you would have a more
		// comprehensive mapping of device types and ports to struct types
		switch r.URL.Path {
		case "/encode/smartlabel/v1":
			switch rawReq.Port {
			case 128:
				var payload smartlabelEncoder.Port128Payload
				if err := json.Unmarshal(rawReq.Payload, &payload); err != nil {
					logger.Logger.Error("error unmarshaling payload", zap.Error(err))
					setBody(w, http.StatusBadRequest, map[string]any{
						"error": fmt.Sprintf("Error unmarshaling payload: %v", err),
						"docs":  "https://docs.truvami.com",
					})
					return
				}
				structPayload = payload
			default:
				logger.Logger.Error("unsupported port", zap.Uint8("port", rawReq.Port))
				setBody(w, http.StatusBadRequest, map[string]any{
					"error": fmt.Sprintf("Unsupported port: %d", rawReq.Port),
					"docs":  "https://docs.truvami.com",
				})
				return
			}
		case "/encode/tagsl/v1":
			switch rawReq.Port {
			case 128:
				var payload tagslEncoder.Port128Payload
				if err := json.Unmarshal(rawReq.Payload, &payload); err != nil {
					logger.Logger.Error("error unmarshaling payload", zap.Error(err))
					setBody(w, http.StatusBadRequest, map[string]any{
						"error": fmt.Sprintf("Error unmarshaling payload: %v", err),
						"docs":  "https://docs.truvami.com",
					})
					return
				}
				structPayload = payload
			case 129:
				var payload tagslEncoder.Port129Payload
				if err := json.Unmarshal(rawReq.Payload, &payload); err != nil {
					logger.Logger.Error("error unmarshaling payload", zap.Error(err))
					setBody(w, http.StatusBadRequest, map[string]any{
						"error": fmt.Sprintf("Error unmarshaling payload: %v", err),
						"docs":  "https://docs.truvami.com",
					})
					return
				}
				structPayload = payload
			case 131:
				var payload tagslEncoder.Port131Payload
				if err := json.Unmarshal(rawReq.Payload, &payload); err != nil {
					logger.Logger.Error("error unmarshaling payload", zap.Error(err))
					setBody(w, http.StatusBadRequest, map[string]any{
						"error": fmt.Sprintf("Error unmarshaling payload: %v", err),
						"docs":  "https://docs.truvami.com",
					})
					return
				}
				structPayload = payload
			case 134:
				var payload tagslEncoder.Port134Payload
				if err := json.Unmarshal(rawReq.Payload, &payload); err != nil {
					logger.Logger.Error("error unmarshaling payload", zap.Error(err))
					setBody(w, http.StatusBadRequest, map[string]any{
						"error": fmt.Sprintf("Error unmarshaling payload: %v", err),
						"docs":  "https://docs.truvami.com",
					})
					return
				}
				structPayload = payload
			default:
				logger.Logger.Error("unsupported port", zap.Uint8("port", rawReq.Port))
				setBody(w, http.StatusBadRequest, map[string]any{
					"error": fmt.Sprintf("Unsupported port: %d", rawReq.Port),
					"docs":  "https://docs.truvami.com",
				})
				return
			}
		default:
			// For other device types, you would add similar switch statements
			logger.Logger.Error("unsupported device type", zap.String("path", r.URL.Path))
			setBody(w, http.StatusBadRequest, map[string]any{
				"error": "Unsupported device type",
				"docs":  "https://docs.truvami.com",
			})
			return
		}

		logger.Logger.Debug("encoding payload", zap.Any("payload", structPayload), zap.Uint8("port", rawReq.Port))

		var warnings []string = nil
		encoded, err := encoder.Encode(structPayload, rawReq.Port)
		if err != nil {
			if errors.Is(err, helpers.ErrValidationFailed) {
				warnings = []string{}
				for _, err := range helpers.UnwrapError(err) {
					logger.Logger.Warn("validation error", zap.Error(err))
					warnings = append(warnings, err.Error())
				}
				logger.Logger.Warn("validation for some fields failed - are you using the correct port?")
			} else {
				logger.Logger.Error("error while encoding payload", zap.Error(err))

				setBody(w, http.StatusBadRequest, map[string]any{
					"error": err.Error(),
					"docs":  "https://docs.truvami.com",
				})
				return
			}
		}

		setBody(w, http.StatusOK, map[string]any{
			"encoded":  encoded,
			"warnings": warnings,
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func setHeaders(w http.ResponseWriter, status int) {
	if status >= 400 {
		w.Header().Set("Content-Type", "application/problem+json")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.WriteHeader(status)
}

func setBody(w http.ResponseWriter, status int, body map[string]any) {
	logger.Logger.Debug("encoding response")

	// add traceId
	traceId := uuid.New().String()
	body["traceId"] = traceId

	data, err := json.Marshal(body)
	if err != nil {
		logger.Logger.Error("error while encoding response", zap.Error(err))
		setHeaders(w, http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))

		if err != nil {
			logger.Logger.Error("error while sending response", zap.Error(err))
		}
		return
	}

	setHeaders(w, status)
	_, err = w.Write(data)
	if err != nil {
		logger.Logger.Error("error while sending response", zap.Error(err))
		return
	}

	logger.Logger.Debug("response sent", zap.Any("response", string(data)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		logger.Logger.Error("error while sending response", zap.Error(err))
	}
}

func loggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// generate a unique request ID
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)

		// start timer
		start := time.Now()

		// use a ResponseWriter wrapper to capture the status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// process the request
		next.ServeHTTP(rw, r)

		// log the details
		logger.Info("HTTP request",
			zap.String("requestId", requestID),
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Int("status", rw.statusCode),
			zap.String("remoteAddress", r.RemoteAddr),
			zap.String("userAgent", r.UserAgent()),
			zap.Duration("latency", time.Since(start)),
		)
	})
}
