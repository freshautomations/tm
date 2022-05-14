package cmd

import (
	"github.com/spf13/cobra"
	"tm/tm/v2/context"
	"tm/tm/v2/startstop"
)

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "Start one or more node(s) or testnet(s)",
	Run: func(cmd *cobra.Command, args []string) {

		// Load chain config
		ctx := context.New(args)

		// Execute start
		startstop.Start(ctx)
	},
}
