package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"tm/m/v2/config"
	"tm/m/v2/ux"
	version "tm/m/v2/version"
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
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(configureCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func LoadConfig() config.Config {
	cfgFile := config.FindConfigFilename()
	fileInfo, err := os.Stat(cfgFile.Path)
	if os.IsNotExist(err) {
		ux.Fatal("%s: no such file or directory", cfgFile.Path)
	}
	if err != nil {
		ux.Fatal(err.Error())
	}
	if fileInfo.IsDir() {
		ux.Fatal("config is a directory: %s", cfgFile.Path)
	}

	// Config defaults
	cfg := config.Config{
		Binary:   "gaiad",
		Home:     "$HOME/.tm",
		Port:     26600,
		Filename: &cfgFile,
	}
	var bytes []byte
	bytes, err = ioutil.ReadFile(cfgFile.Path)
	if err != nil {
		ux.Fatal("could not read config file %s: %s", cfgFile.Path, err)
	}
	err = cfg.CustomUnmarshal(bytes)
	if err != nil {
		ux.Fatal("could not unmarshal config file %s: %s", cfgFile.Path, err)
	}
	cfg.Validate()
	cfg.SetPorts()
	return cfg
}
