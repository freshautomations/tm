package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"tm/tm/v2/ux"
	version "tm/tm/v2/version"
)

var rootCmd = &cobra.Command{
	Use:     "tm",
	Short:   "Testnets Manager",
	Version: version.Version,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var (
	flagHome   string
	flagConfig string
	flagDebug  bool
	flagQuiet  bool
)

func init() {
	// Global (persistent) flags
	// --debug
	rootCmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "", false, "debug log")
	err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		ux.Fatal("could not bind debug flag")
	}
	// --quiet -q
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "quiet operation")
	err = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	if err != nil {
		ux.Fatal("could not bind quiet flag")
	}

	// --home
	homeDir, _ := os.UserHomeDir()
	rootCmd.PersistentFlags().StringVarP(&flagHome, "home", "", filepath.FromSlash(fmt.Sprintf("%s/.tm", homeDir)), "home directory")
	err = viper.BindPFlag("home", rootCmd.PersistentFlags().Lookup("home"))
	if err != nil {
		ux.Fatal("could not bind home flag")
	}

	// --config
	rootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", filepath.FromSlash(fmt.Sprintf("%s/.tm/config.toml", homeDir)), "config file")
	err = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		ux.Fatal("could not bind config flag")
	}

	// sub-commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
