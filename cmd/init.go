package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/fs"
	"io/ioutil"
	"mvdan.cc/sh/v3/shell"
	"os"
	"os/exec"
	"path/filepath"
	"tm/m/v2/config"
	"tm/m/v2/ux"
)

func init() {
	var err error

	// --home
	homeDir, _ := os.UserHomeDir()
	initCmd.Flags().StringVarP(&flagHome, "home", "", fmt.Sprintf("%s/.tm", homeDir), "home directory")
	err = viper.BindPFlag("home", initCmd.Flags().Lookup("home"))
	if err != nil {
		ux.Fatal("could not bind home flag")
	}

	// --config
	initCmd.Flags().StringVarP(&flagConfig, "config", "c", fmt.Sprintf("%s/.tm/config.toml", homeDir), "config file")
	err = viper.BindPFlag("config", initCmd.Flags().Lookup("config"))
	if err != nil {
		ux.Fatal("could not bind config flag")
	}
}

func FindBinary(name string) (result string) {
	expanded, err := shell.Expand(name, nil)
	if err != nil {
		ux.Fatal("%s cannot be expanded ", name)
	}
	result, err = exec.LookPath(expanded)
	if err == nil {
		result, _ = filepath.Abs(result)
	}
	if result == "" {
		result = expanded
	}
	return
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new manager configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// Create default config
		cfg := config.Config{
			Binary: FindBinary("gaiad"),
			Chains: map[string]config.ChainConfig{"testnet-1": {
				StopMaintain: true,
				Nodes: map[string]config.Node{"validator1": {
					Validator: true,
				}},
			}, "testnet-2": {
				Nodes: map[string]config.Node{"validator1": {
					Validator: true,
				}, "fullnode1": {}},
			}},
			Hermes: []config.HermesConfig{{
				Binary: FindBinary("hermes"),
				Nodes:  []string{"testnet-1.validator1", "fullnode1"},
			}},
		}
		cfg.Validate()
		cfg.SetPorts()

		// Encode config
		bytes, err := cfg.CustomMarshal()
		if err != nil {
			ux.Fatal("could not encode initial config: %s", err)
		}

		// Write config
		cfgFile := config.FindConfigFilename()

		_, err = os.Stat(cfgFile.Dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(cfgFile.Dir, fs.ModeDir|fs.ModePerm)
			if err != nil {
				ux.Fatal("could not create config file directory %s", cfgFile.Dir)
			}
		} else {
			if err != nil {
				ux.Fatal("could not create config file directory %s", cfgFile.Dir)
			}
		}
		_, err = os.Stat(cfgFile.Path)
		if os.IsNotExist(err) {
			err = ioutil.WriteFile(cfgFile.Path, bytes, fs.ModePerm)
			if err != nil {
				ux.Fatal("could not write config file %s: %s", cfgFile.Path, err.Error())
			}
		}
	},
}
