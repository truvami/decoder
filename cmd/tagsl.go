package cmd

import (
	"log/slog"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func init() {
	rootCmd.AddCommand(tagslCmd)
}

var tagslCmd = &cobra.Command{
	Use:   "tagsl [port] [payload]",
	Short: "decode tag S / L payloads",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("initializing tagsl decoder")
		d := tagsl.NewTagSLv1Decoder(
			tagsl.WithAutoPadding(AutoPadding),
			tagsl.WithSkipValidation(SkipValidation),
		)

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
