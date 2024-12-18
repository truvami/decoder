package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(nomadxsCmd)
}

var nomadxsCmd = &cobra.Command{
	Use:   "nomadxs [port] [payload]",
	Short: "decode nomad XS payloads",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		logger.Logger.Debug("initializing nomadxs decoder")
		d := nomadxs.NewNomadXSv1Decoder(
			nomadxs.WithAutoPadding(AutoPadding),
		)

		port, err := strconv.Atoi(args[0])
		if err != nil {
			logger.Logger.Error("error while parsing port", zap.Error(err), zap.String("port", args[0]))
			return
		}
		logger.Logger.Debug("port parsed successfully", zap.Int("port", port))

		data, metadata, err := d.Decode(args[1], int16(port), "")
		if err != nil {
			logger.Logger.Error("error while decoding data", zap.Error(err))
			return
		}

		printJSON(data, metadata)
	},
}
