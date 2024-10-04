package cmd

import (
	"log/slog"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/loracloud"
)

var accessToken string

func init() {
	tagxlCmd.Flags().StringVar(&accessToken, "token", "", "Access token for the loracloud API")
	rootCmd.AddCommand(tagxlCmd)
}

var tagxlCmd = &cobra.Command{
	Use:   "tagxl [port] [payload] [devEui] --token [token]",
	Short: "decode tag XL payloads",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))
		if err != nil {
			slog.Error("error while binding token flag", slog.Any("error", err))
			return
		}

		slog.Debug("initilaizing tagxl decoder")
		d := tagxl.NewTagXLv1Decoder(
			loracloud.NewLoracloudMiddleware(accessToken),
		)

		port, err := strconv.Atoi(args[0])
		if err != nil {
			slog.Error("error while parsing port", slog.Any("error", err), slog.String("port", args[0]))
			return
		}
		slog.Debug("port parsed successfully", slog.Int("port", port))

		data, metadata, err := d.Decode(args[1], int16(port), args[2])
		if err != nil {
			slog.Error("error while decoding data", slog.Any("error", err))
			return
		}

		printJSON(data, metadata)
	},
}
