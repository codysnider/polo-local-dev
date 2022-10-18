package config

import (
	"github.com/poloniex/polo-local-dev/config"
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "config",
	Short: "Config management",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var reload = &cobra.Command{
	Use:   "reload",
	Short: "Reload config from dist",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		output.Title("Config")
		output.Plain("Reloading from dist")

		config.Reload()
	},
}

func init() {
	Command.AddCommand(reload)
}
