package cmd

import (
	"github.com/spf13/cobra"
	"tm/tm/v2/context"
	"tm/tm/v2/startstop"
)

var resetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"unsafe-reset", "unsafe-reset-all"},
	Short:   "Reset one or more node(s) or testnet(s) database",
	Run: func(cmd *cobra.Command, args []string) {

		// Load chain config
		ctx := context.New(args)

		// Execute reset
		startstop.Reset(ctx)
	},
}
