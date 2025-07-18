package cmd

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/truvami/decoder/internal/logger"
	helpers "github.com/truvami/decoder/pkg/common"
	"github.com/truvami/decoder/pkg/decoder"
	tagxl "github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/solver"
	"github.com/truvami/decoder/pkg/solver/aws"
	"github.com/truvami/decoder/pkg/solver/loracloud"
	"go.uber.org/zap"
)

var tagXlDevEui string

func init() {
	tagxlCmd.Flags().StringVar(&tagXlDevEui, "dev-eui", "", "DevEUI of the originator device.\nThis is only required for loracloud solver.")
	rootCmd.AddCommand(tagxlCmd)
}

var tagxlCmd = &cobra.Command{
	Use:   "tagxl [port] [payload]",
	Short: "decode tag XL payloads",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
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
			solver, err = loracloud.NewLoracloudClient(ctx, LoracloudAccessToken, logger.Logger)
			if err != nil {
				logger.Logger.Error("error while creating LoRa Cloud position estimate client", zap.Error(err))
				os.Exit(1)
			}
		}

		logger.Logger.Debug("initializing tagxl decoder")
		d := tagxl.NewTagXLv1Decoder(ctx, solver, logger.Logger, tagxl.WithSkipValidation(SkipValidation))

		port, err := strconv.Atoi(args[0])
		if err != nil {
			logger.Logger.Error("error while parsing port", zap.Error(err), zap.String("port", args[0]))
			return
		}
		logger.Logger.Debug("port parsed successfully", zap.Int("port", port))
		if port < 0 || port > 255 {
			logger.Logger.Error("port must be between 0 and 255", zap.Int("port", port))
			return
		}

		ctx = context.WithValue(ctx, decoder.DEVEUI_CONTEXT_KEY, tagXlDevEui)
		ctx = context.WithValue(ctx, decoder.PORT_CONTEXT_KEY, uint8(port))
		ctx = context.WithValue(ctx, decoder.FCNT_CONTEXT_KEY, 1) // Default frame count, can be adjusted as needed

		data, err := d.Decode(ctx, args[1], uint8(port))
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

		printJSON(data.Data)
	},
}
