package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/internal/selfupdate"
	"go.uber.org/zap"
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
var SkipValidation bool

var Solver string
var LoracloudAccessToken string

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

	rootCmd.PersistentFlags().BoolVarP(&SkipValidation, "skip-validation", "", false, "Skip length validation of payload. (default: \033[31mfalse\033[0m)")
	err = viper.BindPFlag("skip-validation", rootCmd.PersistentFlags().Lookup("skip-validation"))
	if err != nil {
		logger.Logger.Error("error while binding skip-validation flag", zap.Error(err))
	}

	rootCmd.PersistentFlags().StringVarP(&Solver, "solver", "s", "aws", "Solver to use for decoding the payload.\nThis can be aws or loracloud.")
	err = viper.BindPFlag("solver", rootCmd.PersistentFlags().Lookup("solver"))
	if err != nil {
		logger.Logger.Error("error while binding solver flag", zap.Error(err))
	}

	rootCmd.PersistentFlags().StringVarP(&LoracloudAccessToken, "loracloud-access-token", "", "", "Loracloud access token to use for decoding the payload. (default: \033[31mempty\033[0m)")
	err = viper.BindPFlag("loracloud-access-token", rootCmd.PersistentFlags().Lookup("loracloud-access-token"))
	if err != nil {
		logger.Logger.Error("error while binding loracloud-access-token flag", zap.Error(err))
	}
}

var rootCmd = &cobra.Command{
	Use:     "decoder",
	Short:   "truvami payload decoder cli helper",
	Version: Version,
	Long: getBanner() + `

A CLI tool to help decode @truvami payloads.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		options := []logger.Option{}

		if Debug {
			options = append(options, logger.WithDebug())
		}

		logger.NewLogger(options...)

		// Non-blocking update check (ignore network errors).
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			latest, has, err := selfupdate.CheckForUpdate(ctx, Version, false)
			if err != nil {
				return
			}
			if has {
				logger.Logger.Info("a new version is available",
					zap.String("current", Version),
					zap.String("latest", latest),
					zap.String("hint", "run 'decoder update' to upgrade"),
				)
			}
		}()

		defer logger.Sync()
	},
}

func Execute() {
	// Initialize logger first to ensure it's available for all commands
	logger.NewLogger()
	defer logger.Sync()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func printJSON(data any) {
	payload := map[string]any{"data": data}

	if Json {
		marshaled, err := json.Marshal(payload)
		if err != nil {
			fmt.Fprintln(os.Stderr, "marshaling error")
			os.Exit(1)
		}
		fmt.Println(string(marshaled))
		return
	}

	marshaled, err := json.MarshalIndent(payload, "", "   ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "marshaling error")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(string(marshaled))
	fmt.Println()
}

func getBanner() string {
	if time.Now().Month() == time.December {
		banner = []string{
			"",
			"\033[1;31m                          ___\033[0m",
			"\033[1;31m                        /`   `'\\\033[0m",
			"\033[1;32m   _                   \033[1;31m/   _..---;      \033[1;32m_\033[0m",
			"\033[1;32m  | |_ _ __ _   ___   _\033[1;31m|  /\033[1;0m__..._/\033[1;32m ___ (_)\033[0m",
			"\033[1;32m  | __| '__| | | \\ \\ / \033[1;31m|.'\033[1;32m| |  _   _ \\| |\033[0m",
			"\033[1;32m  | |_| |  | |_| |\\ V \033[1;0m(_)\033[1;32m_| | | | | | | |\033[0m",
			"\033[1;32m   \\__|_|   \\__,_| \\_/ \\__,_|_| |_| |_|_|\033[0m",
			"",
		}
	}
	return strings.Join(banner, "\n")
}
