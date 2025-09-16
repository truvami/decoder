package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/truvami/decoder/pkg/solver/aws"
)

var capturedAtString string

func init() {
	gnssDebugCmd.Flags().StringVarP(&capturedAtString, "captured-at", "c", time.Now().Format(time.RFC3339), "The time when the GNSS data was captured. (default: current time)")

	rootCmd.AddCommand(gnssDebugCmd)
}

var gnssDebugCmd = &cobra.Command{
	Use:   "gnss-debug",
	Short: "Debug GNSS payloads (experimental)",
	Long: `Debug GNSS payloads by extracting and displaying the 30-bit words from the payload.
	This command is experimental and may change in the future.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		capturedAt, err := time.Parse(time.RFC3339, capturedAtString)
		if err != nil {
			cmd.PrintErrf("Invalid captured-at time format: %v\n", err)
			return
		}

		// if capturedAt is zero, use current time
		if capturedAt.IsZero() {
			capturedAt = time.Now()
		}

		t, err := aws.SolveCapturedAt([]aws.GNSSCapture{
			{
				HexPayload: args[0],
				ReceivedAt: capturedAt,
			},
		})
		if err != nil {
			cmd.PrintErrf("Error solving GNSS payload: %v\n", err)
			return
		}

		cmd.Printf("Extracted time: %s\n", t.Format(time.RFC3339))
	},
}
