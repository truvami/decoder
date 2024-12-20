package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/truvami/decoder/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var banner = []string{
	"\033[32m  _                                   _ ",
	" | |_ _ __ _   ___   ____ _ _ __ ___ (_)",
	" | __| '__| | | \\ \\ / / _` | '_ ` _ \\| |",
	" | |_| |  | |_| |\\ V / (_| | | | | | | |",
	"  \\__|_|   \\__,_| \\_/ \\__,_|_| |_| |_|_|\033[0m",
}

var Debug bool
var Json bool
var AutoPadding bool
var SkipValidation bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: \033[31mfalse\033[0m)")
	err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		logger.Logger.Error("error while binding debug flag", zap.Error(err))
	}

	rootCmd.PersistentFlags().BoolVarP(&Json, "json", "j", false, "Output the result in JSON format. (default: \033[31mfalse\033[0m)")
	err = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	if err != nil {
		logger.Logger.Error("error while binding json flag", zap.Error(err))
	}

	rootCmd.PersistentFlags().BoolVarP(&AutoPadding, "auto-padding", "", false, "Enable automatic padding of payload. (default: \033[31mfalse\033[0m)\n\033[33mWarning:\033[0m this may lead to corrupted data.")
	err = viper.BindPFlag("auto-padding", rootCmd.PersistentFlags().Lookup("auto-padding"))
	if err != nil {
		logger.Logger.Error("error while binding auto-padding flag", zap.Error(err))
	}

	rootCmd.PersistentFlags().BoolVarP(&SkipValidation, "skip-validation", "", false, "Skip length validation of payload. (default: false)")
	err = viper.BindPFlag("skip-validation", rootCmd.PersistentFlags().Lookup("skip-validation"))
	if err != nil {
		slog.Error("error while binding skip-validation flag", slog.Any("error", err))
	}
}

var rootCmd = &cobra.Command{
	Use:   "decoder",
	Short: "truvami payload decoder cli helper",
	Long: strings.Join(banner, "\n") + `

A CLI tool to help decode @truvami payloads.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		options := []logger.Option{}

		if Debug {
			options = append(options, logger.WithDebug())
		}

		if Json {
			// create a custom encoder
			encoderConfig := zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "", // disable stack traces
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
			}

			options = append(options, logger.WithEncoder(zapcore.NewJSONEncoder(encoderConfig)))
		}

		logger.NewLogger(options...)
		defer logger.Sync()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Logger.Error("error while executing command", zap.Error(err))
		os.Exit(1)
	}
}

func printJSON(data interface{}, metadata interface{}) {
	if Json {
		logger.Logger.Info("successfully decoded payload", zap.Any("data", data), zap.Any("metadata", metadata))
		return
	}

	logger.Logger.Info("successfully decoded payload")

	// print data and metadata beautifully and formatted
	marshaled, err := json.MarshalIndent(map[string]interface{}{
		"data":     data,
		"metadata": metadata,
	}, "", "   ")

	// handle marshaling error
	if err != nil {
		logger.Logger.Fatal("marshaling error", zap.Error(err))
	}

	// print the marshaled data
	fmt.Println()
	fmt.Println(string(marshaled))
	fmt.Println()
}
