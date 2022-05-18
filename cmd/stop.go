package cmd

import (
	"github.com/spf13/cobra"
	"tm/tm/v2/context"
	"tm/tm/v2/startstop"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"stop"},
	Short:   "Stop one or more node(s) or testnet(s)",
	Run: func(cmd *cobra.Command, args []string) {

		// Load chain config
		ctx := context.New(args)

		// Execute stop
		startstop.Stop(ctx)
	},
}
