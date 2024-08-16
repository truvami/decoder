package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/truvami/decoder/pkg/logger"
)

var Debug bool
var Verbose bool
var Json bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().BoolVarP(&Json, "json", "j", false, "Output the result in JSON format. (default: false)")
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
}

var rootCmd = &cobra.Command{
	Use:   "decoder",
	Short: "truvami payload decoder cli helper",
	Long:  `A CLI tool to help decode truvami payloads.`,
}

func Execute() {
	cobra.OnInitialize(func() {
		opts := slog.HandlerOptions{
			Level: slog.LevelInfo,
		}

		if Debug {
			opts.Level = slog.LevelDebug
			opts.AddSource = true
		}

		var handler slog.Handler
		if Json {
			handler = slog.NewJSONHandler(os.Stdout, &opts)
		} else {
			handler = logger.NewHandler(&opts)
		}

		slog.SetDefault(slog.New(handler))
	})

	if err := rootCmd.Execute(); err != nil {
		slog.Error("error while executing command", slog.Any("error", err))
		os.Exit(1)
	}
}

func printJSON(data interface{}) {
	slog.Info("successfully decoded payload", slog.Any("data", data))
}
