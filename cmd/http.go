package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/truvami/decoder/internal/logger"
	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	nomadxl "github.com/truvami/decoder/pkg/decoder/nomadxl/v1"
	nomadxs "github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
	smartlabel "github.com/truvami/decoder/pkg/decoder/smartlabel/v1"
	tagsl "github.com/truvami/decoder/pkg/decoder/tagsl/v1"
	tagxl "github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/loracloud"
	"go.uber.org/zap"
)

var host string
var port uint16
var health bool

func init() {
	httpCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the HTTP server to")
	httpCmd.Flags().Uint16Var(&port, "port", 8080, "Port to bind the HTTP server to")
	httpCmd.Flags().StringVar(&accessToken, "token", "", "Access token for the loracloud API")
	httpCmd.Flags().BoolVar(&health, "health", false, "Enable /health endpoint")
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start the HTTP server for the decoder.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(accessToken) == 0 {
			logger.Logger.Warn("no access token provided for loracloud API")
		}

		router := http.NewServeMux()

		// health endpoint
		if health {
			router.HandleFunc("/health", healthHandler)
		}

		type decoderEndpoint struct {
			path    string
			decoder decoder.Decoder
		}

		var decoders []decoderEndpoint = []decoderEndpoint{
			{
				"tagxl/v1",
				tagxl.NewTagXLv1Decoder(
					loracloud.NewLoracloudMiddleware(accessToken),
					tagxl.WithAutoPadding(AutoPadding),
					tagxl.WithSkipValidation(SkipValidation),
				),
			},
			{"tagsl/v1", tagsl.NewTagSLv1Decoder(
				tagsl.WithAutoPadding(AutoPadding),
				tagsl.WithSkipValidation(SkipValidation),
			)},
			{"nomadxs/v1", nomadxs.NewNomadXSv1Decoder(
				nomadxs.WithAutoPadding(AutoPadding),
				nomadxs.WithSkipValidation(SkipValidation),
			)},
			{"nomadxl/v1", nomadxl.NewNomadXLv1Decoder(
				nomadxl.WithAutoPadding(AutoPadding),
				nomadxl.WithSkipValidation(SkipValidation),
			)},
			{"smartlabel/v1", smartlabel.NewSmartLabelv1Decoder(
				loracloud.NewLoracloudMiddleware(accessToken),
				smartlabel.WithAutoPadding(AutoPadding),
				smartlabel.WithSkipValidation(SkipValidation),
			)},
		}

		// add the decoders
		for _, d := range decoders {
			addDecoder(router, d.path, d.decoder)
		}

		// middleware
		handler := loggingMiddleware(logger.Logger, router)

		logger.Logger.Info("starting HTTP server", zap.String("host", host), zap.Uint64("port", uint64(port)))
		err := http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), handler)

		if err != nil {
			logger.Logger.Error("error while starting HTTP server", zap.Error(err))
			os.Exit(1)
		}
	},
}

func addDecoder(router *http.ServeMux, path string, decoder decoder.Decoder) {
	logger.Logger.Debug("adding decoder", zap.String("path", path))
	router.HandleFunc("POST /"+path, getHandler(decoder))
}

func getHandler(decoder decoder.Decoder) func(http.ResponseWriter, *http.Request) {
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

		logger.Logger.Debug("decoding payload")

		var warnings []string = nil
		data, err := decoder.Decode(req.Payload, req.Port, req.DevEUI)
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
			"metadata": data.Metadata,
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
