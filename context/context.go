package context

import (
	"fmt"
	"strings"
	"tm/tm/v2/config"
	"tm/tm/v2/utils"
	"tm/tm/v2/ux"
)

type Context struct {
	AllNodeNames      []string
	AllValidatorNames []string
	AllChainNames     []string
	Input             []string
	Config            config.Config
}

// New creates a new context and loads the configuration
func New(args []string) Context {
	// Check if there were any arguments specified.
	allNodesArg := len(args) == 0
	// Trim input
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	// Create context defaults
	ctx := Context{
		Config: config.New(),
	}
	// Get all nodes, validators and chain names
	for chainName, chain := range ctx.Config.Chains {
		ctx.AllChainNames = append(ctx.AllChainNames, chainName)
		for nodeName := range chain.Nodes {
			ctx.AllNodeNames = append(ctx.AllNodeNames, fmt.Sprintf("%s.%s", chainName, nodeName))
			if chain.Nodes[nodeName].Validator {
				ctx.AllValidatorNames = append(ctx.AllValidatorNames, fmt.Sprintf("%s.%s", chainName, nodeName))
			}
		}
	}
	// If no arguments were specified, run the command for all nodes.
	if allNodesArg {
		ctx.Input = ctx.AllNodeNames
	}
	// Fill in Input based on the input.
	for _, arg := range args {
		// Input is a testnet name, add all nodes
		if utils.Contains(ctx.AllChainNames, arg) {
			for nodeName := range ctx.Config.Chains[arg].Nodes {
				ctx.Input = append(ctx.Input, fmt.Sprintf("%s.%s", arg, nodeName))
			}
			continue
		}
		fullArgName, err := utils.FindNodeFullname(ctx.AllNodeNames, arg)
		if err != nil {
			ux.Fatal("invalid input %s", arg)
		}
		ctx.Input = append(ctx.Input, fullArgName)
	}
	// Make sure input is unique
	var inputResult []string
	for _, input := range ctx.Input {
		if !utils.Contains(inputResult, input) {
			inputResult = append(inputResult, input)
		}
	}
	ctx.Input = inputResult
	return ctx
}
