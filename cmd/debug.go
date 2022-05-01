package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tm/m/v2/config"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug command for development",
	Run: func(cmd *cobra.Command, args []string) {

		// Create debug config
		cfg := config.Config{
			Binary:       FindBinary("gaiad"),
			Home:         "/Users/greg/.tm",
			StopMaintain: false,
			Wallets: []config.Wallet{{
				Name:      "wallet1",
				Mnemonics: "q w e",
			}, {
				Name:      "wallet2",
				Mnemonics: "g h j",
			}},
			Chains: map[string]config.ChainConfig{"testnet-1": {
				HDPath: "myhdpath",
				Binary: "gaiad",
				Home:   "t1home",
				Nodes: map[string]config.Node{"validator1": {
					Validator:   true,
					Binary:      "custombianry",
					Home:        "val1home",
					Mnemonics:   "a b c",
					Port:        26000,
					Connections: []string{"fullnode2", "fullnode1"},
				}, "fullnode1": {
					Validator:    false,
					Binary:       "gaiad",
					Home:         "fullhome",
					StopMaintain: true,
					Mnemonics:    "",
					Port:         26650,
					Connections:  nil,
				}, "fullnode2": {
					Connections: []string{"testnet-1.fullnode1"},
				}},
			}, "testnet-2": {
				HDPath:       "myhdpath3",
				Binary:       "regen3",
				Home:         "home3",
				StopMaintain: true,
				Nodes: map[string]config.Node{"validator1": {
					Validator: true,
				}, "fullnode1": {}},
			}},
			Hermes: []config.HermesConfig{{
				Binary:           FindBinary("hermes"),
				Config:           "$HOME/.hermes/config.toml",
				LogLevel:         "info",
				TelemetryEnabled: true,
				TelemetryHost:    "localhost",
				TelemetryPort:    3001,
				Nodes:            []string{"testnet-2.validator1", "testnet-1.fullnode1"},
			}, {
				Binary:           FindBinary("hermesx"),
				Config:           "$HOME/.hermes2/config.toml",
				LogLevel:         "debug",
				TelemetryEnabled: false,
				TelemetryHost:    "localhost",
				TelemetryPort:    3002,
				Nodes:            []string{"testnet-1.validator1", "testnet-2.fullnode1"},
			}},
			Port:     26600,
			Filename: &config.Filename{},
		}
		cfg.Validate()
		cfg.SetPorts()

		bytes, err := cfg.CustomMarshal()
		if err != nil {
			panic(err)
		}

		var cfg2 config.Config
		err = cfg2.CustomUnmarshal(bytes)
		if err != nil {
			panic(err)
		}
		bytes, err = cfg2.CustomMarshal()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))
	},
}
