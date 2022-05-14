package cmd

import (
	"github.com/spf13/cobra"
	"tm/tm/v2/ux"
	"tm/tm/v2/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		ux.Info(version.Version)
	},
}
