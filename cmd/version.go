package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tm/m/v2/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
