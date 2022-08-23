package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var JsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "pld",
	Short: "Poloniex Local Dev Toolkit",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVar(&JsonOutput, "json", false, "JSON output")
}
