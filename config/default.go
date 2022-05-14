package config

import (
	"tm/tm/v2/tmconfig"
	"tm/tm/v2/utils"
)

func NewDefaultConfig() Config {
	filename := tmconfig.FindConfigFilename()
	cfg := Config{
		Binary: utils.FindOSBinary("gaiad"),
		Chains: map[string]*ChainConfig{"testnet-1": {
			StopMaintain: true,
			Nodes: map[string]*Node{"validator1": {
				Validator: true,
			}},
		}, "testnet-2": {
			Nodes: map[string]*Node{"validator1": {
				Validator: true,
			}, "fullnode1": {}},
		}},
		Hermes: []HermesConfig{{
			Binary: utils.FindOSBinary("hermes"),
			Nodes:  []string{"testnet-1.validator1", "fullnode1"},
		}},
		Filename: &filename,
	}
	cfg.validate()
	cfg.setPorts()

	return cfg
}
