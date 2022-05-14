package config

import (
	"fmt"
	"mvdan.cc/sh/v3/shell"
	"strings"
	"tm/m/v2/utils"
	"tm/m/v2/ux"
)

// validate checks the node logic in the configuration file.
func (cfg *Config) validate() {

	// Trim extra spaces from strings
	var allChains []string
	cfg.Binary = strings.TrimSpace(cfg.Binary)
	cfg.Home = strings.TrimSpace(cfg.Home)
	for _, wallet := range cfg.Wallets {
		wallet.Name = strings.TrimSpace(wallet.Name)
		wallet.Mnemonics = strings.TrimSpace(wallet.Mnemonics)
	}
	for chainName, chain := range cfg.Chains {
		chain.HDPath = strings.TrimSpace(chain.HDPath)
		chain.Binary = strings.TrimSpace(chain.Binary)
		chain.Home = strings.TrimSpace(chain.Home)
		chain.denom = strings.TrimSpace(chain.denom)
		allChains = append(allChains, chainName)
		for nodeName, node := range chain.Nodes {
			node.Binary = strings.TrimSpace(node.Binary)
			node.Home = strings.TrimSpace(node.Home)
			node.Mnemonics = strings.TrimSpace(node.Mnemonics)
			if node.Port > 65535 {
				ux.Fatal("invalid port %s in chain %s node %s config", node.Port, chainName, nodeName)
			}
			for i := range node.Connections {
				node.Connections[i] = strings.TrimSpace(node.Connections[i])
			}
		}
	}
	for _, hermes := range cfg.Hermes {
		hermes.Binary = strings.TrimSpace(hermes.Binary)
		hermes.Config = strings.TrimSpace(hermes.Config)
		hermes.LogLevel = strings.TrimSpace(hermes.LogLevel)
		hermes.TelemetryHost = strings.TrimSpace(hermes.TelemetryHost)
		hermes.Mnemonics = strings.TrimSpace(hermes.Mnemonics)
		if hermes.TelemetryPort > 65535 {
			ux.Fatal("invalid port %s in hermes config", hermes.TelemetryPort)
		}
		for i := range hermes.Nodes {
			hermes.Nodes[i] = strings.TrimSpace(hermes.Nodes[i])
		}
	}
	if cfg.Port > 65535 {
		ux.Fatal("invalid port %s in global config", cfg.Port)
	}

	// Wallet names are unique
	// Wallet mnemonics are unique
	allWallets := make([]string, 0)
	allWalletMnemonics := make([]string, 0)
	for _, wallet := range cfg.Wallets {
		if utils.Contains(allWallets, wallet.Name) {
			ux.Fatal("duplicate wallet name %s", wallet)
		}

		if wallet.Mnemonics != "" && utils.Contains(allWalletMnemonics, wallet.Mnemonics) {
			ux.Fatal("duplicate wallet mnemonic for wallet %s", wallet)
		}

		allWallets = append(allWallets, wallet.Name)
		allWalletMnemonics = append(allWalletMnemonics, wallet.Mnemonics)
	}

	// Chain IDs are inherently unique, no validation necessary.
	// Node names are unique within a chain.
	// Node names do not match chain IDs.
	// There is at least one validator per chain.
	var allNodes []string
	for chainID, chain := range cfg.Chains {
		foundValidator := false
		for moniker, node := range chain.Nodes {
			if utils.Contains(allChains, moniker) {
				ux.Fatal("chain name and node name cannot both match %s", moniker)
			}
			nodeFullname := fmt.Sprintf("%s.%s", chainID, moniker)
			if utils.Contains(allNodes, nodeFullname) {
				ux.Fatal("duplicate node moniker %s.%s", chainID, moniker)
			}
			allNodes = append(allNodes, nodeFullname)
			if node.Validator {
				foundValidator = true
			}
		}

		if !foundValidator {
			ux.Fatal("at least one validator required at %s definition", chainID)
		}
	}

	// Node connections have to point to other valid nodes within the same chain.
	// Node connections might not mention their chain ID since all connections are within the same chain.
	// Node connections do not point to self.
	// Only one node connection to one server. (no repeat)
	for chainID, chain := range cfg.Chains {
		for nodeMoniker, node := range chain.Nodes {
			var connections []string
			for i, connection := range node.Connections {
				if len(strings.Split(connection, ".")) == 1 {
					connection = fmt.Sprintf("%s.%s", chainID, connection)
				}
				connectionFullname, err := utils.FindNodeFullname(allNodes, connection)
				if err != nil {
					ux.Fatal("%s connection %s", nodeMoniker, err.Error())
				}
				if connectionFullname == fmt.Sprintf("%s.%s", chainID, nodeMoniker) {
					ux.Fatal("connection %s in node %s.%s points to self", connection, chainID, nodeMoniker)
				}
				if !utils.Contains(allNodes, connectionFullname) {
					ux.Fatal("non-existent connection %s", connection)
				}
				connectionFullnameSplit := strings.Split(connectionFullname, ".")
				if connectionFullnameSplit[0] != chainID {
					ux.Fatal("connection %s in node %s.%s points to other network", connection, chainID, nodeMoniker)
				}
				if utils.Contains(connections, connectionFullname) {
					ux.Fatal("connection %s is duplicated in node %s.%s", connection, chainID, nodeMoniker)
				}
				connections = append(connections, connectionFullname)
				node.Connections[i] = connectionFullnameSplit[1]
			}
		}
	}

	// Each Hermes config should have at least one node
	// Hermes Config parameter is unique
	// Hermes nodes connect to valid nodes only
	// Hermes points to maximum one node per chain
	var allHermesConfig []string
	for i, hermes := range cfg.Hermes {

		if len(hermes.Nodes) == 0 {
			ux.Fatal("no Hermes nodes at %d.[[Hermes]] definition", i+1)
		}

		expanded, err := shell.Expand(hermes.Config, nil)
		if err != nil {
			ux.Fatal("config cannot be expanded at %d.[[Hermes]] definition", i+1)
		}
		if utils.Contains(allHermesConfig, expanded) {
			ux.Fatal("config path has to be unique at %d.[[Hermes]] definition", i+1)
		}
		allHermesConfig = append(allHermesConfig, expanded)

		var allHermesNetworks []string
		var connectionFullname string
		for j, connection := range hermes.Nodes {
			connectionFullname, err = utils.FindNodeFullname(allNodes, connection)
			if err != nil {
				ux.Fatal("%s at %d.[[Hermes]] definition", err.Error(), i+1)
			}
			hermes.Nodes[j] = connectionFullname
			if !utils.Contains(allNodes, connectionFullname) {
				ux.Fatal("non-existent connection %s at %d.[[Hermes]] definition", connection, i+1)
			}
			connectionFullnameSplit := strings.Split(connectionFullname, ".")
			if len(connectionFullnameSplit) != 2 {
				ux.Fatal("invalid connection name at %d.[[Hermes]] definition", connection, i+1)
			}
			connectionChainID := connectionFullnameSplit[0]
			if utils.Contains(allHermesNetworks, connectionChainID) {
				ux.Fatal("multiple node connection to %s chain at %d.[[Hermes]] definition", connectionChainID, i+1)
			}
			allHermesNetworks = append(allHermesNetworks, connectionChainID)
		}
	}
}
