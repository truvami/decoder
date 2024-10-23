package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/truvami/decoder/pkg/logger"
)

var banner = []string{
	"  _                                   _ ",
	" | |_ _ __ _   ___   ____ _ _ __ ___ (_)",
	" | __| '__| | | \\ \\ / / _` | '_ ` _ \\| |",
	" | |_| |  | |_| |\\ V / (_| | | | | | | |",
	"  \\__|_|   \\__,_| \\_/ \\__,_|_| |_| |_|_|",
}

var Debug bool
var Verbose bool
var Json bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		slog.Error("error while binding verbose flag", slog.Any("error", err))
	}

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	err = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		slog.Error("error while binding debug flag", slog.Any("error", err))
	}

	rootCmd.PersistentFlags().BoolVarP(&Json, "json", "j", false, "Output the result in JSON format. (default: false)")
	err = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	if err != nil {
		slog.Error("error while binding json flag", slog.Any("error", err))
	}
}

var rootCmd = &cobra.Command{
	Use:   "decoder",
	Short: "truvami payload decoder cli helper",
	Long: strings.Join(banner, "\n") + `

A CLI tool to help decode @truvami payloads.`,
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

func printJSON(data interface{}, metadata interface{}) {
	slog.Info("successfully decoded payload", slog.Any("data", data), slog.Any("metadata", metadata))
}
