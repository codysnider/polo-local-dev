package start

import (
	"github.com/poloniex/polo-local-dev/output"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "start",
	Short: "Start project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		output.Title("Start")

	},
}
