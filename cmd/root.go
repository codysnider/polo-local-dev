package cmd

import (
	"github.com/poloniex/polo-local-dev/cmd/build"
	"github.com/poloniex/polo-local-dev/cmd/clone"
	"github.com/poloniex/polo-local-dev/cmd/config"
	"github.com/poloniex/polo-local-dev/cmd/dependency"
	"github.com/poloniex/polo-local-dev/cmd/doctor"
	"github.com/poloniex/polo-local-dev/cmd/fork"
	"github.com/poloniex/polo-local-dev/cmd/project"
	"github.com/poloniex/polo-local-dev/cmd/start"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "pld",
		Short: "Poloniex Local Dev Toolkit",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Config
	rootCmd.AddCommand(config.Command)

	// Doctor
	rootCmd.AddCommand(doctor.Command)

	// Fork
	rootCmd.AddCommand(fork.Command)

	// Clone
	rootCmd.AddCommand(clone.Command)

	// Start
	rootCmd.AddCommand(start.Command)

	// Build
	rootCmd.AddCommand(build.Command)

	// Dependency
	rootCmd.AddCommand(dependency.Command)

	// Project
	rootCmd.AddCommand(project.Command)

	// Global flags
	rootCmd.PersistentFlags().BoolP("json", "j", false, "JSON output")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
