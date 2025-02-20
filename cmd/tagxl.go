package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/truvami/decoder/internal/logger"
	helpers "github.com/truvami/decoder/pkg/common"
	tagxl "github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/loracloud"
	"go.uber.org/zap"
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
			logger.Logger.Error("error while binding token flag", zap.Error(err))
			return
		}

		logger.Logger.Debug("initializing tagxl decoder")
		d := tagxl.NewTagXLv1Decoder(
			loracloud.NewLoracloudMiddleware(accessToken),
			tagxl.WithAutoPadding(AutoPadding),
			tagxl.WithSkipValidation(SkipValidation),
		)

		port, err := strconv.Atoi(args[0])
		if err != nil {
			logger.Logger.Error("error while parsing port", zap.Error(err), zap.String("port", args[0]))
			return
		}
		logger.Logger.Debug("port parsed successfully", zap.Int("port", port))

		data, err := d.Decode(args[1], int16(port), args[2])
		if err != nil {
			if errors.Is(err, helpers.ErrValidationFailed) {
				for _, err := range helpers.UnwrapError(err) {
					logger.Logger.Warn("", zap.Error(err))
				}
				logger.Logger.Warn("validation for some fields failed - are you using the correct port?")
			} else {
				logger.Logger.Error("error while decoding data", zap.Error(err))
				return
			}
		}

		printJSON(data.Data, data.Metadata)
	},
}
