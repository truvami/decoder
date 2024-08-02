package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "decoder",
	Short: "truvami payload decoder cli helper",
	Long:  `A CLI tool to help decode truvami payloads.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printJSON(data interface{}) {
	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
