package cmd

import (
	"github.com/spf13/cobra"
	"tm/tm/v2/config"
	"tm/tm/v2/context"
	"tm/tm/v2/initialize"
	"tm/tm/v2/tmconfig"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize"},
	Short:   "Initialize new manager configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// Create new tm default config, if necessary
		tmconfig.CreateConfigPath()
		cfg := config.NewDefaultConfig()
		cfg.SaveNotOverwrite()

		// Load chain config
		ctx := context.New(args)

		// Initialize chain config
		initialize.Initialize(ctx)
	},
}
