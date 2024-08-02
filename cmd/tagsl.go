package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func init() {
	rootCmd.AddCommand(tagslCmd)
}

var tagslCmd = &cobra.Command{
	Use:   "tagsl [port] [payload]",
	Short: "decode tag S / L payloads",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		d := tagsl.NewTagSLv1Decoder()

		port, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("error parsing port: %v", err)
			return
		}

		data, err := d.Decode(args[1], int16(port), "")
		if err != nil {
			fmt.Printf("error decoding data: %v", err)
			return
		}

		printJSON(data)
	},
}
