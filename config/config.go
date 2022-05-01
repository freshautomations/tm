package config

import (
	"fmt"
	"mvdan.cc/sh/v3/shell"
	"os/exec"
	"path/filepath"
	"strings"
	"tm/m/v2/ux"
)

// Config defines the Testnets Manager configuration.
type Config struct {
	Binary           string                 `toml:"binary,omitempty"`
	Home             string                 `toml:"home,omitempty"`
	StopMaintain     bool                   `toml:"stop_maintain,omitempty"`
	NoConfigOverride bool                   `toml:"no_config_override,omitempty"`
	Wallets          []Wallet               `toml:"wallet,omitempty"`
	Chains           map[string]ChainConfig `toml:"-"`
	Hermes           []HermesConfig         `toml:"hermes,omitempty"`
	Port             uint                   `toml:"port,omitzero"`
	Filename         *Filename              `toml:"-"`
}

// HermesConfig defines the Hermes-related entries in the configuration file.
type HermesConfig struct {
	Binary           string   `toml:"binary,omitempty"`
	Config           string   `toml:"config,omitempty"`
	LogLevel         string   `toml:"log_level,omitempty"`
	TelemetryEnabled bool     `toml:"telemetry_enabled,omitempty"`
	TelemetryHost    string   `toml:"telemetry_host,omitempty"`
	TelemetryPort    uint     `toml:"telemetry_port,omitzero"`
	Nodes            []string `toml:"nodes,omitempty"`
}

// ChainConfig defines the Testnets Manager chain configuration format
type ChainConfig struct {
	HDPath       string          `toml:"hdpath,omitempty"`
	Binary       string          `toml:"binary,omitempty"`
	Home         string          `toml:"home,omitempty"`
	StopMaintain bool            `toml:"stop_maintain,omitempty"`
	Nodes        map[string]Node `toml:"-"`
}

type Node struct {
	Validator    bool     `toml:"validator,omitempty"`
	Binary       string   `toml:"binary,omitempty"`
	Home         string   `toml:"home,omitempty"`
	StopMaintain bool     `toml:"stop_maintain,omitempty"`
	Mnemonics    string   `toml:"mnemonics,omitempty"` // Not used on full nodes
	Port         uint     `toml:"port,omitzero"`
	Connections  []string `toml:"connections,omitempty"` // default is to connect all validators to each other and all full nodes to all validators
}

type Wallet struct {
	Name      string `toml:"name"`
	Mnemonics string `toml:"mnemonics,omitempty"`
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func findNodeFullname(allNodes []string, item string) (string, error) {
	if item == "" {
		return "", fmt.Errorf("empty item not allowed")
	}
	itemSplit := strings.Split(item, ".")
	if len(itemSplit) > 2 {
		return "", fmt.Errorf("too many fields in item %s", item)
	}
	if len(itemSplit) == 2 {
		if contains(allNodes, item) {
			return item, nil
		} else {
			return "", fmt.Errorf("node name not found %s", item)
		}
	}
	// itemSplit == 1, item contains the node name (without chain ID)
	found := false
	result := ""
	for _, fullName := range allNodes {
		nameSplit := strings.Split(fullName, ".")
		if len(nameSplit) != 2 {
			return "", fmt.Errorf("invalid node name %s", fullName)
		}
		if nameSplit[1] == item {
			if found {
				return "", fmt.Errorf("ambivalent node moniker %s", item)
			} else {
				found = true
				result = fullName
			}
		}
	}
	if found {
		return result, nil
	}
	return "", fmt.Errorf("node name not found %s", item)
}

func (cfg Config) findNodeItem(fullNodeName string) (ChainConfig, Node) {
	fullNodeNameSplit := strings.Split(fullNodeName, ".")
	if len(fullNodeNameSplit) != 2 {
		ux.Fatal("invalid node name %s", fullNodeName)
	}
	chainName := fullNodeNameSplit[0]
	nodeNameClean := fullNodeNameSplit[1]

	for chainNameLoop, chain := range cfg.Chains {
		if chainNameLoop == chainName {
			for nodeNameLoop, node := range chain.Nodes {
				if nodeNameLoop == nodeNameClean {
					return chain, node
				}
			}
		}
	}
	ux.Fatal("node name %s not found in config", fullNodeName)
	return ChainConfig{}, Node{}
}

// Validate checks the node logic in the configuration file.
func (cfg Config) Validate() {
	var allNodes []string

	// Trim extra spaces from strings
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
	var allWallets []string
	var allWalletMnemonics []string
	for _, wallet := range cfg.Wallets {

		if contains(allWallets, wallet.Name) {
			ux.Fatal("duplicate wallet name %s", wallet)
		}

		if wallet.Mnemonics != "" && contains(allWalletMnemonics, wallet.Mnemonics) {
			ux.Fatal("duplicate wallet mnemonic for wallet %s", wallet)
		}

		allWallets = append(allWallets, wallet.Name)
		allWalletMnemonics = append(allWalletMnemonics, wallet.Mnemonics)
	}

	// Chain IDs are inherently unique, no validation necessary.
	// Node names are unique within a chain.
	// There is at least one validator per chain.
	for chainID, chain := range cfg.Chains {
		foundValidator := false
		for moniker, node := range chain.Nodes {
			nodeFullname := fmt.Sprintf("%s.%s", chainID, moniker)
			if contains(allNodes, nodeFullname) {
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
	// Node connections may not mention their chain ID since all connections are within the same chain.
	// Node connections do not point to self.
	// Only one node connection to one server. (no repeat)
	for chainID, chain := range cfg.Chains {
		for nodeMoniker, node := range chain.Nodes {
			var connections []string
			for i, connection := range node.Connections {
				if len(strings.Split(connection, ".")) == 1 {
					connection = fmt.Sprintf("%s.%s", chainID, connection)
				}
				connectionFullname, err := findNodeFullname(allNodes, connection)
				if err != nil {
					ux.Fatal("%s connection %s", nodeMoniker, err.Error())
				}
				if connectionFullname == fmt.Sprintf("%s.%s", chainID, nodeMoniker) {
					ux.Fatal("connection %s in node %s.%s points to self", connection, chainID, nodeMoniker)
				}
				if !contains(allNodes, connectionFullname) {
					ux.Fatal("non-existent connection %s", connection)
				}
				connectionFullnameSplit := strings.Split(connectionFullname, ".")
				if connectionFullnameSplit[0] != chainID {
					ux.Fatal("connection %s in node %s.%s points to other network", connection, chainID, nodeMoniker)
				}
				if contains(connections, connectionFullname) {
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
		if contains(allHermesConfig, expanded) {
			ux.Fatal("config path has to be unique at %d.[[Hermes]] definition", i+1)
		}
		allHermesConfig = append(allHermesConfig, expanded)

		var allHermesNetworks []string
		var connectionFullname string
		for j, connection := range hermes.Nodes {
			connectionFullname, err = findNodeFullname(allNodes, connection)
			if err != nil {
				ux.Fatal("%s at %d.[[Hermes]] definition", err.Error(), i+1)
			}
			hermes.Nodes[j] = connectionFullname
			if !contains(allNodes, connectionFullname) {
				ux.Fatal("non-existent connection %s at %d.[[Hermes]] definition", connection, i+1)
			}
			connectionFullnameSplit := strings.Split(connectionFullname, ".")
			if len(connectionFullnameSplit) != 2 {
				ux.Fatal("invalid connection name at %d.[[Hermes]] definition", connection, i+1)
			}
			connectionChainID := connectionFullnameSplit[0]
			if contains(allHermesNetworks, connectionChainID) {
				ux.Fatal("multiple node connection to %s chain at %d.[[Hermes]] definition", connectionChainID, i+1)
			}
			allHermesNetworks = append(allHermesNetworks, connectionChainID)
		}
	}
}

func (cfg *Config) SetPorts() {
	// Set up ports
	if cfg.Port == 0 {
		cfg.Port = 26600
	}
	initialCfgPort := cfg.Port
	for chainName, chain := range cfg.Chains {
		for nodeName, node := range chain.Nodes {
			if node.Port == 0 {
				node.Port = cfg.Port
				cfg.Port += 10
			} else {
				if node.Port > initialCfgPort && node.Port < cfg.Port {
					ux.Fatal("Preset port conflicts with automatically assigned ports between (%d-%d) at %s.%s", initialCfgPort, cfg.Port, chainName, nodeName)
				}
			}
			chain.Nodes[nodeName] = node
		}
		cfg.Chains[chainName] = chain
	}
}

func (cfg Config) GetHome(nodeFullName string) string {
	chain, node := cfg.findNodeItem(nodeFullName)
	result := node.Home
	if result == "" {
		result = chain.Home
	}
	if result == "" {
		result = cfg.Home
	}
	if result == "" {
		result = "$HOME/.tm" // default
	}
	return result
}

func (cfg Config) GetBinary(nodeFullName string) string {
	chain, node := cfg.findNodeItem(nodeFullName)
	result := node.Binary
	if result == "" {
		result = chain.Binary
	}
	if result == "" {
		result = cfg.Binary
	}
	if result == "" {
		var err error
		result, err = exec.LookPath("gaiad")
		if err == nil {
			result, _ = filepath.Abs(result)
		}
		if result == "" {
			result = "gaiad"
		}
	}
	return result
}

func (cfg Config) GetStopMaintain(nodeFullName string) bool {
	chain, node := cfg.findNodeItem(nodeFullName)
	result := node.StopMaintain
	if result == false {
		result = chain.StopMaintain
	}
	if result == false {
		result = cfg.StopMaintain
	}
	return result
}
