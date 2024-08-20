package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	sloghttp "github.com/samber/slog-http"
	"github.com/spf13/cobra"
	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
	"github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/loracloud"
)

var host string
var port uint16

func init() {
	httpCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the HTTP server to")
	httpCmd.Flags().Uint16Var(&port, "port", 8080, "Port to bind the HTTP server to")
	httpCmd.Flags().StringVar(&accessToken, "token", "", "Access token for the loracloud API")
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start the HTTP server for the decoder.",
	Run: func(cmd *cobra.Command, args []string) {
		router := http.NewServeMux()

		type decoderEndpoint struct {
			path    string
			decoder decoder.Decoder
		}

		var decoders []decoderEndpoint = []decoderEndpoint{
			{
				"tagxl/v1",
				tagxl.NewTagXLv1Decoder(
					loracloud.NewLoracloudMiddleware(accessToken),
				),
			},
			{"tagsl/v1", tagsl.NewTagSLv1Decoder()},
			{"nomadxs/v1", nomadxs.NewNomadXSv1Decoder()},
		}

		// add the decoders
		for _, d := range decoders {
			slog.Debug("adding decoder", slog.String("path", d.path))
			addDecoder(router, d.path, d.decoder)
		}

		// middleware
		handler := sloghttp.Recovery(router)
		handler = sloghttp.New(slog.Default())(handler)

		slog.Info("starting HTTP server", slog.String("host", host), slog.Uint64("port", uint64(port)))
		err := http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), handler)

		if err != nil {
			slog.Error("error while starting HTTP server", slog.Any("error", err))
			os.Exit(1)
		}
	},
}

func addDecoder(router *http.ServeMux, path string, decoder decoder.Decoder) {
	slog.Debug("adding decoder", slog.String("path", path))
	router.HandleFunc("POST /"+path, getHandler(decoder))
}

func getHandler(decoder decoder.Decoder) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Port    int16  `json:"port"`
			Payload string `json:"payload"`
			DevEUI  string `json:"devEui"`
		}

		// decode the request
		var req request

		slog.Debug("decoding request")
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error("error while decoding request", slog.Any("error", err))
			setHeaders(w, http.StatusBadRequest)
			return
		}

		// decode the payload
		slog.Debug("decoding payload")
		data, err := decoder.Decode(req.Payload, req.Port, req.DevEUI)
		if err != nil {
			slog.Error("error while decoding payload", slog.Any("error", err))
			setHeaders(w, http.StatusBadRequest)
			return
		}

		// data to json
		slog.Debug("encoding response")
		data, err = json.Marshal(data)
		if err != nil {
			slog.Error("error while encoding response", slog.Any("error", err))
			setHeaders(w, http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// send the response
		setHeaders(w, http.StatusOK)
		w.Write(data.([]byte))

		slog.Debug("response sent", slog.Any("response", string(data.([]byte))))
	}
}

func setHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.WriteHeader(status)
}
