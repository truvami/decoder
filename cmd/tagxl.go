package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/loracloud"
)

var accessToken string

func init() {
	tagxlCmd.PersistentFlags().StringVar(&accessToken, "token", "", "Access token for the loracloud API")
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	rootCmd.AddCommand(tagxlCmd)
}

var tagxlCmd = &cobra.Command{
	Use:   "tagxl [port] [payload] [devEui] --token [token]",
	Short: "decode tag XL payloads",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		d := tagxl.NewTagXLv1Decoder(
			loracloud.NewLoracloudMiddleware(accessToken),
		)

		port, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("error parsing port: %v", err)
			return
		}

		data, err := d.Decode(args[1], int16(port), args[2])
		if err != nil {
			fmt.Printf("error decoding data: %v", err)
			return
		}

		printJSON(data)
	},
}
