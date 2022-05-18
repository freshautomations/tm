package cmd

import (
	"github.com/spf13/cobra"
	"tm/tm/v2/context"
	"tm/tm/v2/startstop"
)

var (
	flagf bool
	flagF bool
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Get the log of one or more node(s) or testnet(s)",
	Run: func(cmd *cobra.Command, args []string) {

		// Load chain config
		ctx := context.New(args)

		// Execute log
		startstop.Log(ctx)
	},
}
