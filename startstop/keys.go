package startstop

import (
	"encoding/json"
	"strings"
	"tm/tm/v2/context"
	"tm/tm/v2/execute"
	"tm/tm/v2/utils"
	"tm/tm/v2/ux"
)

func Keys(ctx context.Context) {
	var doneNetworkNames []string
	for _, fullNodename := range ctx.Input {
		fullNodenameSplit := strings.Split(fullNodename, ".")
		chainName := fullNodenameSplit[0]
		if utils.Contains(doneNetworkNames, chainName) {
			continue
		}
		result, err := execute.KeysList(ctx.Config.GetBinary(fullNodename), ctx.Config.GetChainHome(fullNodename))
		if err != nil {
			ux.Warn("could not retrieve keys from %s: %s", chainName, err)
		}
		var data []utils.Keys
		err = json.Unmarshal([]byte(result), &data)
		if err != nil {
			ux.Warn("could not unmarshal keys output from %s: %s", chainName, err)
			ux.Warn("Output string: %s", result)
		}
		for _, d := range data {
			ux.Info("- name: %s", d.Name)
			ux.Info("  type: %s", d.KeyType)
			ux.Info("  address: %s", d.Address)
			ux.Info("  pubkey: '%s'", d.PubKey)
			ux.Info("  mnemonic: \"%s\"", ctx.Config.GetMnemonics(chainName, d.Name))
		}
		doneNetworkNames = append(doneNetworkNames, chainName)
	}
}
