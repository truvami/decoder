package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/truvami/decoder/internal/logger"
	"github.com/truvami/decoder/internal/selfupdate"
	"go.uber.org/zap"
)

var (
	updateNext bool
)

func init() {
	updateCmd.Flags().BoolVar(&updateNext, "next", false, "Include release candidates (e.g. -rc) when checking/updating.")
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update decoder to the latest available version",
	Long: `Update decoder to the latest GitHub release for your OS/ARCH.

By default, only stable releases are considered. Pass --next to include release candidates (tags with "-rc").`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// First check for update (stable by default, or include RCs with --next)
		latest, has, err := selfupdate.CheckForUpdate(ctx, Version, updateNext)
		if err != nil {
			logger.Logger.Error("failed to check for updates", zap.Error(err))
			return
		}
		if !has {
			logger.Logger.Info("already up to date",
				zap.String("current", Version),
				zap.String("latest", latest),
			)
			return
		}

		// Increase timeouts for the actual download and replacement
		selfupdate.SetClientTimeout(2 * time.Minute)
		ctxUpdate, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		newTag, err := selfupdate.UpdateToLatest(ctxUpdate, Version, updateNext)
		if err != nil {
			// This may include a Windows-specific message where the new binary is placed next to the existing one.
			logger.Logger.Error("update failed", zap.Error(err))
			return
		}

		logger.Logger.Info("update completed",
			zap.String("from", Version),
			zap.String("to", newTag),
			zap.String("note", "restart any running instance to use the new version"),
		)
	},
}
