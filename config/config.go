package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mvdan.cc/sh/v3/shell"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"tm/tm/v2/consts"
	"tm/tm/v2/tmconfig"
	"tm/tm/v2/utils"
	"tm/tm/v2/ux"
)

// Config defines the Testnets Manager configuration.
type Config struct {
	Binary           string                  `toml:"binary,omitempty"`
	Home             string                  `toml:"home,omitempty"`
	StopMaintain     bool                    `toml:"stop_maintain,omitempty"`
	NoConfigOverride bool                    `toml:"no_config_override,omitempty"`
	Wallets          []Wallet                `toml:"wallet,omitempty"`
	Chains           map[string]*ChainConfig `toml:"-"`
	Hermes           []HermesConfig          `toml:"hermes,omitempty"`
	Port             uint                    `toml:"port,omitzero"`
	Filename         *tmconfig.Filename      `toml:"-"`
}

// HermesConfig defines the Hermes-related entries in the configuration file.
type HermesConfig struct {
	Binary           string   `toml:"binary,omitempty"`
	Config           string   `toml:"config,omitempty"`
	LogLevel         string   `toml:"log_level,omitempty"`
	TelemetryEnabled bool     `toml:"telemetry_enabled,omitempty"`
	TelemetryHost    string   `toml:"telemetry_host,omitempty"`
	TelemetryPort    uint     `toml:"telemetry_port,omitzero"`
	Mnemonics        string   `toml:"mnemonics,omitempty"`
	Nodes            []string `toml:"nodes,omitempty"`
}

// ChainConfig defines the Testnets Manager chain configuration format
type ChainConfig struct {
	HDPath       string           `toml:"hdpath,omitempty"`
	Binary       string           `toml:"binary,omitempty"`
	Home         string           `toml:"home,omitempty"`
	StopMaintain bool             `toml:"stop_maintain,omitempty"`
	denom        string           `toml:"denom,omitempty"`
	Nodes        map[string]*Node `toml:"-"`
}

type Node struct {
	Binary       string   `toml:"binary,omitempty"`
	Home         string   `toml:"home,omitempty"`
	Mnemonics    string   `toml:"mnemonics,omitempty"` // Not used on full nodes
	Port         uint     `toml:"port,omitzero"`
	Validator    bool     `toml:"validator,omitempty"`
	StopMaintain bool     `toml:"stop_maintain,omitempty"`
	Connections  []string `toml:"connections,omitempty"` // default is to connect all validators to each other and all full nodes to all validators
}

type Wallet struct {
	Name      string `toml:"name"`
	Mnemonics string `toml:"mnemonics,omitempty"`
}

func New() Config {
	cfgFile := tmconfig.FindConfigFilename()
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
	cfg := Config{
		Binary:   "gaiad",
		Home:     tmconfig.FindConfigFilename().Dir,
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
	cfg.validate()
	cfg.setPorts()
	return cfg
}

func (cfg *Config) setPorts() {
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

func (cfg Config) FindNode(fullNodeName string) (*ChainConfig, *Node) {
	fullNodeNameSplit := strings.Split(fullNodeName, ".")
	if len(fullNodeNameSplit) != 2 {
		ux.Fatal("invalid node name %s", fullNodeName)
	}
	chainName := fullNodeNameSplit[0]
	nodeName := fullNodeNameSplit[1]

	chain, ok := cfg.Chains[chainName]
	if !ok {
		ux.Fatal("chain for node %s not found in config", fullNodeName)
	}
	node, ok2 := cfg.Chains[chainName].Nodes[nodeName]
	if !ok2 {
		ux.Fatal("node %s not found in config", fullNodeName)
	}
	return chain, node
}

func (cfg Config) GetHome(nodeFullName string) string {
	chain, node := cfg.FindNode(nodeFullName)
	nodeFullnameSplit := strings.Split(nodeFullName, ".")
	chainName := nodeFullnameSplit[0]
	nodeName := nodeFullnameSplit[1]
	result := ""
	if node.Home != "" {
		result = utils.GetSlashPath(node.Home)
	} else {
		if chain.Home != "" {
			result = utils.GetSlashPath("%s/%s", chain.Home, nodeName)
		} else {
			if cfg.Home != "" {
				result = utils.GetSlashPath("%s/%s/%s", cfg.Home, chainName, nodeName)
			} else {
				result = utils.GetSlashPath("%s/%s/%s", tmconfig.FindConfigFilename().Dir, chainName, nodeName)
			}
		}
	}
	expanded, err := shell.Expand(result, nil)
	if err != nil {
		ux.Fatal(err.Error())
	}
	return expanded
}

// GetChainHome returns the home folder for the chain. Input can be "ChainName" or "ChainName.NodeName" format.
func (cfg Config) GetChainHome(nodeFullName string) string {
	nodeFullnameSplit := strings.Split(nodeFullName, ".")
	chainName := nodeFullnameSplit[0]
	chain := cfg.Chains[chainName]
	result := ""
	if chain.Home != "" {
		result = chain.Home
	} else {
		if cfg.Home != "" {
			result = utils.GetSlashPath("%s/%s", cfg.Home, chainName)
		} else {
			result = utils.GetSlashPath("%s/%s", tmconfig.FindConfigFilename().Dir, chainName)
		}
	}
	expanded, err := shell.Expand(result, nil)
	if err != nil {
		ux.Fatal(err.Error())
	}
	return expanded
}

func (cfg Config) GetBinary(nodeFullName string) string {
	chain, node := cfg.FindNode(nodeFullName)
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
	var err error
	result, err = shell.Expand(result, nil)
	if err != nil {
		ux.Fatal("binary not found, %s", err.Error())
	}
	return result
}

// GetChainBinary returns the binary associated with the chain. Input can be "ChainName" or "ChainName.NodeName" format.
func (cfg Config) GetChainBinary(nodeFullName string) string {
	nodeFullNameSplit := strings.Split(nodeFullName, ".")
	chainName := nodeFullNameSplit[0]
	chain := cfg.Chains[chainName]
	result := chain.Binary
	if result == "" {
		result = cfg.Binary
	}
	if result == "" {
		if len(nodeFullNameSplit) > 1 {
			nodeName := nodeFullNameSplit[1]
			result = chain.Nodes[nodeName].Binary
		}
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
	var err error
	result, err = shell.Expand(result, nil)
	if err != nil {
		ux.Fatal("binary not found, %s", err.Error())
	}
	return result
}

func (cfg Config) GetStopMaintain(nodeFullName string) bool {
	chain, node := cfg.FindNode(nodeFullName)
	result := node.StopMaintain
	if result == false {
		result = chain.StopMaintain
	}
	if result == false {
		result = cfg.StopMaintain
	}
	return result
}

func (cfg Config) GetPort(nodeFullName string) uint {
	_, node := cfg.FindNode(nodeFullName)
	return node.Port
}

func (cfg Config) GetRPCPort(nodeFullName string) uint {
	return cfg.GetPort(nodeFullName)
}

func (cfg Config) GetAppPort(nodeFullName string) uint {
	return cfg.GetPort(nodeFullName) + 1
}

func (cfg Config) GetGRPCPort(nodeFullName string) uint {
	return cfg.GetPort(nodeFullName) + 2
}

func (cfg Config) GetP2PPort(nodeFullName string) uint {
	return cfg.GetPort(nodeFullName) + 3
}

func (cfg Config) GetPPROFPort(nodeFullName string) uint {
	return cfg.GetPort(nodeFullName) + 4
}

func (cfg Config) GetKMSPort(nodeFullName string) uint {
	return cfg.GetPort(nodeFullName) + 5
}

func (cfg Config) GetGRPCWEBPort(nodeFullName string) uint {
	return cfg.GetPort(nodeFullName) + 6
}

func (cfg Config) GetPath(fullNodename string, suffix string) string {
	return utils.GetSlashPath("%s/%s", cfg.GetHome(fullNodename), suffix)
}

func (cfg Config) GetChainPath(fullNodename string, suffix string) string {
	return utils.GetSlashPath("%s/%s", cfg.GetChainHome(fullNodename), suffix)
}

func (cfg Config) GetDenom(fullNodename string) string {
	chainName := strings.Split(fullNodename, ".")[0]
	chain := cfg.Chains[chainName]
	if chain.denom != "" {
		return chain.denom
	}
	chainGenesis := cfg.GetChainPath(fullNodename, "config/genesis.json")
	if denom, ok := utils.GetConfigEntry(chainGenesis, "app_state.staking.params.bond_denom").(string); !ok {
		ux.Fatal("cannot get denomination in genesis")
	} else {
		return denom
	}
	ux.Fatal("invalid code segment during GetDenom")
	return ""
}

func (cfg Config) GetConnections(fullNodename string) []string {
	chain, node := cfg.FindNode(fullNodename)

	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]
	nodeName := fullNodenameSplit[1]

	if len(node.Connections) > 0 {
		return node.Connections
	}

	// If no connection was specified, connect to all validator nodes, except self.
	var result []string
	for nodeNameLoop, nodeLoop := range chain.Nodes {
		if nodeNameLoop == nodeName {
			continue
		}
		if nodeLoop.Validator {
			result = append(result, fmt.Sprintf("%s.%s", chainName, nodeNameLoop))
		}
	}
	return result
}

func (cfg Config) GetMnemonics(chainName string, walletName string) string {
	home := cfg.GetChainHome(chainName)
	file, err := ioutil.ReadFile(consts.GetMnemonics(home, walletName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		ux.Fatal("could not read mnemonics %s: %s", walletName, err)
	}

	var data utils.Keys
	err = json.Unmarshal(file, &data)
	if err != nil {
		ux.Warn("could not unmarshal keys output from %s: %s", chainName, err)
		ux.Warn("Output string: %s", string(file))
	}
	return data.Mnemonic
}
