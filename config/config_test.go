package config

import (
	"testing"
	"tm/tm/v2/tmconfig"
	"tm/tm/v2/utils"
	"tm/tm/v2/ux"
)

func newDebugConfig() Config {
	cfg := Config{
		Binary:       utils.FindOSBinary("gaiad"),
		Home:         utils.GetSlashPath("$HOME/.tm"),
		StopMaintain: false,
		Wallets: []Wallet{{
			Name:      "wallet1",
			Mnemonics: "q w e",
		}, {
			Name:      "wallet2",
			Mnemonics: "g h j",
		}},
		Chains: map[string]*ChainConfig{"testnet-1": {
			HDPath: "myhdpath",
			Binary: "gaiad",
			Home:   "t1home",
			Nodes: map[string]*Node{"validator1": {
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
			Nodes: map[string]*Node{"validator1": {
				Validator: true,
			}, "fullnode1": {}},
		}},
		Hermes: []HermesConfig{{
			Binary:           utils.FindOSBinary("hermes"),
			Config:           "$HOME/.hermes/config.toml",
			LogLevel:         "info",
			TelemetryEnabled: true,
			TelemetryHost:    "localhost",
			TelemetryPort:    3001,
			Nodes:            []string{"testnet-2.validator1", "testnet-1.fullnode1"},
		}, {
			Binary:           utils.FindOSBinary("hermesx"),
			Config:           "$HOME/.hermes2/config.toml",
			LogLevel:         "debug",
			TelemetryEnabled: false,
			TelemetryHost:    "localhost",
			TelemetryPort:    3002,
			Nodes:            []string{"testnet-1.validator1", "testnet-2.fullnode1"},
		}},
		Port:     26600,
		Filename: &tmconfig.Filename{},
	}
	cfg.validate()
	cfg.setPorts()
	return cfg
}

func TestConfig(t *testing.T) {
	// Create debug config
	cfg := newDebugConfig()
	cfg.validate()
	cfg.setPorts()

	bytes, err := cfg.CustomMarshal()
	if err != nil {
		panic(err)
	}

	var cfg2 Config
	err = cfg2.CustomUnmarshal(bytes)
	if err != nil {
		panic(err)
	}
	bytes, err = cfg2.CustomMarshal()
	if err != nil {
		panic(err)
	}
	ux.Info(string(bytes))
}
