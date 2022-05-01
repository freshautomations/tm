package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure testnets configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := LoadConfig()
		fmt.Printf("%+v", cfg)
		fmt.Printf("configuration:\n%v\n", cfg.Chains["testnet-1"].Nodes)
	},
}
