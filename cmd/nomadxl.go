package cmd

import (
	"log/slog"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/truvami/decoder/pkg/decoder/nomadxl/v1"
)

func init() {
	rootCmd.AddCommand(nomadxlCmd)
}

var nomadxlCmd = &cobra.Command{
	Use:   "nomadxl [port] [payload]",
	Short: "decode nomad XL payloads",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("initilaizing nomadxs decoder")
		d := nomadxl.NewNomadXLv1Decoder()

		port, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error("error while parsing port", slog.Any("error", err), slog.String("port", args[0]))
			return
		}
		slog.Debug("port parsed successfully", slog.Int("port", port))

		data, metadata, err := d.Decode(args[1], int16(port), "")
		if err != nil {
			slog.Error("error while decoding data", slog.Any("error", err))
			return
		}

		printJSON(data, metadata)
	},
}