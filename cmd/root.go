package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"tm/tm/v2/utils"
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
	rootCmd.PersistentFlags().StringVarP(&flagHome, "home", "", utils.GetSlashPath("%s/.tm", homeDir), "home directory")
	err = viper.BindPFlag("home", rootCmd.PersistentFlags().Lookup("home"))
	if err != nil {
		ux.Fatal("could not bind home flag")
	}

	// --config
	rootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", utils.GetSlashPath("%s/.tm/config.toml", homeDir), "config file")
	err = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		ux.Fatal("could not bind config flag")
	}

	// -f for log
	logCmd.Flags().BoolVarP(&flagf, "follow", "f", false, "output appended data as the log file grows")
	err = viper.BindPFlag("follow", logCmd.Flags().Lookup("follow"))
	if err != nil {
		ux.Fatal("could not bind follow flag")
	}

	// -F for log
	logCmd.Flags().BoolVarP(&flagF, "follow-and-retry", "F", false, "output appended data as the log file grows and keep trying to open the file if it is inaccessible")
	err = viper.BindPFlag("follow-and-retry", logCmd.Flags().Lookup("follow-and-retry"))
	if err != nil {
		ux.Fatal("could not bind follow-and-retry flag")
	}

	// sub-commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(keysCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
