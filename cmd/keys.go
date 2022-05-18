package cmd

import (
	"github.com/spf13/cobra"
	"tm/tm/v2/context"
	"tm/tm/v2/startstop"
)

var keysCmd = &cobra.Command{
	Use:     "keys",
	Aliases: []string{"k"},
	Short:   "List keys of one or more node(s) or testnet(s)",
	Run: func(cmd *cobra.Command, args []string) {

		// Load chain config
		ctx := context.New(args)

		// Execute start
		startstop.Keys(ctx)
	},
}
