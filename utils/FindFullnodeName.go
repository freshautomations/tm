package utils

import (
	"fmt"
	"strings"
)

// FindNodeFullname looks for a nodename in a list of all nodes and returns its full designation.
// Full node names are presented in a "ChainName.NodeName" format. This function requests a string slice of all nodes
// in this format and a node name to be searched for in the list. The node name can be in the full format or it can
// omit the ChainName. If the ChainName is omitted, the function will only return a full node if the node name is unique
// in the list.
func FindNodeFullname(allNodes []string, item string) (string, error) {
	if item == "" {
		return "", fmt.Errorf("empty item not allowed")
	}
	itemSplit := strings.Split(item, ".")
	if len(itemSplit) > 2 {
		return "", fmt.Errorf("too many fields in item %s", item)
	}
	if len(itemSplit) == 2 {
		if Contains(allNodes, item) {
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
