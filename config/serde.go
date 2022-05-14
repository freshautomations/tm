package config

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"strconv"
	"strings"
)

func extractBool(v interface{}) (bool, error) {
	if v == nil {
		return false, nil
	}
	if rawValue, ok := v.(bool); ok {
		return rawValue, nil
	} else {
		if rawValue, ok := v.(string); !ok {
			return false, fmt.Errorf("could not extract value from %v", v)
		} else {
			value := strings.TrimSpace(strings.ToLower(rawValue))
			if value == "true" || value == "false" {
				return value == "true", nil
			}
			return false, fmt.Errorf("stringified boolean value invalid %v", rawValue)
		}
	}
}

func extractUint(v interface{}) (uint, error) {
	if v == nil {
		return 0, nil
	}
	if rawValueInt64, ok0 := v.(int64); ok0 {
		if rawValueInt64 < 0 {
			return 0, fmt.Errorf("negative values not accepted %d", rawValueInt64)
		}
		return uint(rawValueInt64), nil
	} else {
		if rawValue, ok1 := v.(uint); ok1 {
			return rawValue, nil
		} else {
			if rawValue64, ok2 := v.(uint64); ok2 {
				if rawValue64 > uint64(^uint(0)) {
					return 0, fmt.Errorf("value out of bounds %d", rawValue64)
				}
				return uint(rawValue64), nil
			} else {
				if rawValueFloat64, ok3 := v.(float64); ok3 {
					return uint(rawValueFloat64), nil
				} else {
					if rawValueString, ok4 := v.(string); !ok4 {
						return 0, fmt.Errorf("could not extract value from %v", v)
					} else {
						valueString := strings.TrimSpace(strings.ToLower(rawValueString))
						valueUint64, err := strconv.ParseUint(valueString, 10, 0)
						if err != nil {
							return 0, err
						}
						if valueUint64 > uint64(^uint(0)) {
							return 0, fmt.Errorf("value out of bounds %d", valueUint64)
						}
						return uint(valueUint64), nil
					}
				}
			}
		}
	}
}

func extractStringSlice(v interface{}) ([]string, error) {
	if v == nil {
		return []string{}, nil
	}
	var result []string
	if rawSlice, ok := v.([]interface{}); ok {
		for _, rawValue := range rawSlice {
			if value, ok2 := rawValue.(string); ok2 {
				result = append(result, strings.TrimSpace(value))
			}
		}
		return result, nil
	} else {
		if rawValueString, ok2 := v.(string); !ok2 {
			return []string{}, fmt.Errorf("could not extract value from %v", v)
		} else {
			return []string{strings.TrimSpace(strings.ToLower(rawValueString))}, nil
		}
	}
}

func extractString(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	if result, ok := v.(string); !ok {
		return "", fmt.Errorf("could not extract value from %v", v)
	} else {
		return result, nil
	}
}

func (cfg Config) CustomMarshal() ([]byte, error) {
	// Encode config
	var buf bytes.Buffer
	var err error

	encoder := toml.NewEncoder(&buf)
	err = encoder.Encode(cfg)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(cfg.Chains)
	if err != nil {
		return nil, err
	}
	for chainName, chain := range cfg.Chains {
		for nodeName, node := range chain.Nodes {
			_, _ = buf.Write([]byte(fmt.Sprintf("\n[%s.%s]\n", chainName, nodeName)))
			err = encoder.Encode(node) // This will not indent the values properly. It's a shortcoming of the toml library used.
			if err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), err
}

func (cfg *Config) CustomUnmarshal(data []byte) error {
	// Decode what we can
	meta, err := toml.Decode(string(data), &cfg)
	if err != nil {
		return err
	}

	// Get the whole data decoded, so we can check for invalid values and unparsed items.
	decoded := make(map[string]interface{})
	_, err = toml.Decode(string(data), &decoded)
	if err != nil {
		return err
	}

	// Find chains data
	chains := make(map[string]*ChainConfig)
	for _, key := range meta.Undecoded() {
		keySplit := strings.Split(key.String(), ".")
		switch len(keySplit) {
		case 0:
			return fmt.Errorf("invalid key %s", key)
		case 1: // whole chain
			chainName := keySplit[0]
			if meta.Type(key.String()) != "Hash" {
				continue // not a [chain] definition but some random key
			}
			if chainItem, ok := decoded[chainName].(map[string]interface{}); ok {
				if _, ok = chains[chainName]; ok {
					return fmt.Errorf("%s duplicate definition", chainName)
				}
				// Decode what we can into ChainConfig manually
				var stopMaintain bool
				var hdpath string
				var binary string
				var home string
				var denom string
				stopMaintain, err = extractBool(chainItem["stop_maintain"])
				if err != nil {
					return err
				}
				hdpath, err = extractString(chainItem["hdpath"])
				if err != nil {
					return err
				}
				binary, err = extractString(chainItem["binary"])
				if err != nil {
					return err
				}
				home, err = extractString(chainItem["home"])
				if err != nil {
					return err
				}
				denom, err = extractString(chainItem["denom"])
				if err != nil {
					return err
				}
				emptyNodes := make(map[string]*Node)
				chains[chainName] = &ChainConfig{
					HDPath:       hdpath,
					Binary:       binary,
					Home:         home,
					StopMaintain: stopMaintain,
					Nodes:        emptyNodes,
					denom:        denom,
				}
			}
		case 2: // one node
			chainName := keySplit[0]
			nodeName := keySplit[1]
			if _, ok := chains[chainName]; !ok {
				return fmt.Errorf("%s.%s chain undefined", chainName, nodeName)
			}
			if _, ok := chains[chainName].Nodes[nodeName]; ok {
				return fmt.Errorf("duplicate node %s.%s", chainName, nodeName)
			}
			if chainItem, ok := decoded[chainName].(map[string]interface{}); ok {
				if nodeItem, ok2 := chainItem[nodeName].(map[string]interface{}); ok2 {
					var validator bool
					var stopMaintain bool
					var port uint
					var connections []string
					var mnemonics string
					var binary string
					var home string
					validator, err = extractBool(nodeItem["validator"])
					if err != nil {
						return err
					}
					stopMaintain, err = extractBool(nodeItem["stop_maintain"])
					if err != nil {
						return err
					}
					port, err = extractUint(nodeItem["port"])
					if err != nil {
						return err
					}
					if port > 65535 {
						return fmt.Errorf("invalid port %d at %s.%s", port, chainName, nodeName)
					}
					connections, err = extractStringSlice(nodeItem["connections"])
					if err != nil {
						return err
					}
					mnemonics, err = extractString(nodeItem["mnemonics"])
					if err != nil {
						return err
					}
					binary, err = extractString(nodeItem["binary"])
					if err != nil {
						return err
					}
					home, err = extractString(nodeItem["home"])
					if err != nil {
						return err
					}
					chains[chainName].Nodes[nodeName] = &Node{
						Validator:    validator,
						Binary:       binary,
						Home:         home,
						StopMaintain: stopMaintain,
						Mnemonics:    mnemonics,
						Port:         port,
						Connections:  connections,
					}
				}
			}
		default:
			continue
		}
	}
	cfg.Chains = chains
	cfg.validate()
	cfg.setPorts()
	return nil
}
