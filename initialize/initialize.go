package initialize

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"tm/tm/v2/consts"
	"tm/tm/v2/context"
	"tm/tm/v2/execute"
	"tm/tm/v2/utils"
	"tm/tm/v2/ux"
)

func Initialize(ctx context.Context) {
	var doneNetworkNames []string
	for _, fullNodename := range ctx.Input {
		fullNodenameSplit := strings.Split(fullNodename, ".")
		chainName := fullNodenameSplit[0]

		if !utils.Contains(doneNetworkNames, chainName) {
			runInit(ctx, fullNodename)
			setDenomInChainGenesis(ctx, fullNodename)
			createWallets(ctx, fullNodename)
			addGenesisAccounts(ctx, fullNodename)
			createGentxTransactions(ctx, fullNodename)
			collectGentxs(ctx, fullNodename)
			ValidateGenesis(ctx, fullNodename)
			removeResidualsFromChainFolder(ctx, fullNodename)
			copyGenesis(ctx, fullNodename)
			configure(ctx, fullNodename)
			doneNetworkNames = append(doneNetworkNames, chainName)
		}
	}
}

func runInit(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]
	for nodeNameLoop, node := range ctx.Config.Chains[chainName].Nodes {
		fullNodenameLoop := fmt.Sprintf("%s.%s", chainName, nodeNameLoop)
		binary := ctx.Config.GetBinary(fullNodenameLoop)
		home := ctx.Config.GetHome(fullNodenameLoop)
		execute.Init(fullNodenameLoop, binary, home)
		if node.Validator {
			chainGenesis := ctx.Config.GetChainPath(fullNodenameLoop, "config/genesis.json")
			_, err := os.Stat(chainGenesis)
			if err != nil {
				var bytesRead []byte
				bytesRead, err = ioutil.ReadFile(ctx.Config.GetPath(fullNodenameLoop, "config/genesis.json"))
				if err != nil {
					ux.Fatal("could not read genesis.json at %s", fullNodenameLoop)
				}
				err = os.Mkdir(ctx.Config.GetChainPath(fullNodenameLoop, "config"), fs.ModeDir|fs.ModePerm)
				if err != nil && !errors.Is(err, os.ErrExist) {
					ux.Fatal("could not create chain config folder at %s", ctx.Config.GetChainPath(fullNodenameLoop, "config"))
				}
				err = ioutil.WriteFile(chainGenesis, bytesRead, fs.ModePerm)
				if err != nil {
					ux.Fatal("cannot write chain genesis, %s", err.Error())
				}
			}
		}
	}
	keysFolder := consts.GetMnemonicsDir(chainName)
	err := os.Mkdir(keysFolder, fs.ModeDir|fs.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		ux.Fatal("could not create chain keys folder at %s", keysFolder)
	}
}

func setDenomInChainGenesis(ctx context.Context, fullNodename string) {
	denom := ctx.Config.GetDenom(fullNodename)
	chainGenesis := ctx.Config.GetChainPath(fullNodename, "config/genesis.json")
	utils.SetConfigEntry(chainGenesis, "app_state.crisis.constant_fee.denom", denom)
	type fee struct {
		Amount string `json:"amount"`
		Denom  string `json:"denom"`
	}
	utils.SetConfigEntry(chainGenesis, "app_state.gov.deposit_params.min_deposit", []fee{{Amount: "10000000", Denom: denom}})
	utils.SetConfigEntry(chainGenesis, "app_state.liquidity.params.pool_creation_fee", []fee{{Amount: "40000000", Denom: denom}})
	utils.SetConfigEntry(chainGenesis, "app_state.mint.params.mint_denom", denom)
	utils.SetConfigEntry(chainGenesis, "app_state.staking.params.bond_denom", denom)
}

func createWallets(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]

	hdpath := ctx.Config.Chains[chainName].HDPath
	// Create keys for all validators
	for nodeNameLoop, nodeLoop := range ctx.Config.Chains[chainName].Nodes {
		if nodeLoop.Validator {
			binary := ctx.Config.GetChainBinary(fmt.Sprintf("%s.%s", chainName, nodeNameLoop))
			home := ctx.Config.GetChainHome(fmt.Sprintf("%s.%s", chainName, nodeNameLoop))
			execute.KeysAdd(binary, home, nodeNameLoop, hdpath, nodeLoop.Mnemonics)
		}
	}
	// Create keys for all Hermes instances
	for i, hermes := range ctx.Config.Hermes {
		binary := ctx.Config.GetChainBinary(chainName)
		home := ctx.Config.GetChainHome(chainName)
		execute.KeysAdd(binary, home, fmt.Sprintf("hermes%d", i), hdpath, hermes.Mnemonics)
	}
	// Create keys for all wallets
	for _, wallet := range ctx.Config.Wallets {
		binary := ctx.Config.GetChainBinary(chainName)
		home := ctx.Config.GetChainHome(chainName)
		execute.KeysAdd(binary, home, wallet.Name, hdpath, wallet.Mnemonics)
	}
}

func addGenesisAccounts(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]

	// Add account for validators
	for nodeNameLoop, nodeLoop := range ctx.Config.Chains[chainName].Nodes {
		if nodeLoop.Validator {
			binary := ctx.Config.GetChainBinary(fmt.Sprintf("%s.%s", chainName, nodeNameLoop))
			home := ctx.Config.GetChainHome(fmt.Sprintf("%s.%s", chainName, nodeNameLoop))
			denom := ctx.Config.GetDenom(chainName)
			execute.AddGenesisAccount(binary, home, nodeNameLoop, fmt.Sprintf("10000000000%s,10000000000samoleans", denom))
		}
	}

	// Create account for all Hermes instances on all initializing networks
	// Note that Hermes instances might have connections to other networks that are not initialized here.
	for i := range ctx.Config.Hermes {
		binary := ctx.Config.GetChainBinary(chainName)
		home := ctx.Config.GetChainHome(chainName)
		denom := ctx.Config.GetDenom(chainName)
		execute.AddGenesisAccount(binary, home, fmt.Sprintf("hermes%d", i), fmt.Sprintf("10000000000%s,10000000000samoleans", denom))
	}

	// Create account for all wallets
	for _, wallet := range ctx.Config.Wallets {
		binary := ctx.Config.GetChainBinary(chainName)
		home := ctx.Config.GetChainHome(chainName)
		denom := ctx.Config.GetDenom(chainName)
		execute.AddGenesisAccount(binary, home, wallet.Name, fmt.Sprintf("10000000000%s,10000000000samoleans", denom))
	}
}

func createGentxTransactions(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]

	for nodeNameLoop, nodeLoop := range ctx.Config.Chains[chainName].Nodes {
		if nodeLoop.Validator {
			binary := ctx.Config.GetChainBinary(fmt.Sprintf("%s.%s", chainName, nodeNameLoop))
			home := ctx.Config.GetChainHome(fmt.Sprintf("%s.%s", chainName, nodeNameLoop))
			denom := ctx.Config.GetDenom(chainName)
			execute.AddGentx(binary, home, chainName, nodeNameLoop, fmt.Sprintf("1000000000%s", denom))
		}
	}
}

func removeResidualsFromChainFolder(ctx context.Context, fullNodename string) {
	_ = os.Remove(ctx.Config.GetChainPath(fullNodename, "config/app.toml"))
	_ = os.Remove(ctx.Config.GetChainPath(fullNodename, "config/config.toml"))
	_ = os.Remove(ctx.Config.GetChainPath(fullNodename, "config/priv_validator_key.json"))
	_ = os.Remove(ctx.Config.GetChainPath(fullNodename, "config/client.toml"))
	_ = os.Remove(ctx.Config.GetChainPath(fullNodename, "config/node_key.json"))
	_ = os.Remove(ctx.Config.GetChainPath(fullNodename, "data/priv_validator_state.json"))
	_ = os.Remove(ctx.Config.GetChainPath(fullNodename, "data"))
}

func collectGentxs(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]

	binary := ctx.Config.GetChainBinary(chainName)
	home := ctx.Config.GetChainHome(chainName)
	execute.CollectGentxs(binary, home)
}

func ValidateGenesis(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]

	binary := ctx.Config.GetChainBinary(chainName)
	home := ctx.Config.GetChainHome(chainName)
	execute.ValidateGenesis(binary, home)
}

func copyGenesis(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]

	chainGenesis := ctx.Config.GetChainPath(chainName, "config/genesis.json")
	for nodeName := range ctx.Config.Chains[chainName].Nodes {
		genesis := ctx.Config.GetPath(fmt.Sprintf("%s.%s", chainName, nodeName), "config/genesis.json")
		data, err := ioutil.ReadFile(chainGenesis)
		if err != nil {
			ux.Fatal("could not read genesis for chain %s", chainName)
		}
		err = ioutil.WriteFile(genesis, data, fs.ModePerm)
		if err != nil {
			ux.Fatal("could not write genesis for %s", fullNodename)
		}
	}
}

func configure(ctx context.Context, fullNodename string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]

	for nodeName := range ctx.Config.Chains[chainName].Nodes {
		fullNodename = fmt.Sprintf("%s.%s", chainName, nodeName)

		// config.toml settings
		configToml := ctx.Config.GetPath(fullNodename, "config/config.toml")
		p2pAddress := fmt.Sprintf("tcp://0.0.0.0:%d", ctx.Config.GetP2PPort(fullNodename))
		rpcAddress := fmt.Sprintf("tcp://0.0.0.0:%d", ctx.Config.GetRPCPort(fullNodename))
		pprofAddress := fmt.Sprintf("0.0.0.0:%d", ctx.Config.GetPPROFPort(fullNodename))
		utils.SetConfigEntry(configToml, "p2p.laddr", p2pAddress)
		utils.SetConfigEntry(configToml, "rpc.laddr", rpcAddress)
		utils.SetConfigEntry(configToml, "rpc.pprof_laddr", pprofAddress)

		// app.toml settings
		appToml := ctx.Config.GetPath(fullNodename, "config/app.toml")
		appAddress := fmt.Sprintf("tcp://0.0.0.0:%d", ctx.Config.GetAppPort(fullNodename))
		grpcAddress := fmt.Sprintf("0.0.0.0:%d", ctx.Config.GetGRPCPort(fullNodename))
		grpcWebAddress := fmt.Sprintf("0.0.0.0:%d", ctx.Config.GetGRPCWEBPort(fullNodename))
		minimumGasPrices := fmt.Sprintf("0%s", ctx.Config.GetDenom(fullNodename))

		utils.SetConfigEntry(appToml, "minimum-gas-prices", minimumGasPrices)
		utils.SetConfigEntry(appToml, "api.address", appAddress)
		utils.SetConfigEntry(appToml, "api.enable", true)
		utils.SetConfigEntry(appToml, "api.swagger", true)
		utils.SetConfigEntry(appToml, "grpc.address", grpcAddress)
		utils.SetConfigEntry(appToml, "grpc-web.address", grpcWebAddress)

		if ctx.Config.GetStopMaintain(fullNodename) {
			continue
		}

		var peers []string
		var peerIDs []string
		for _, fullNodenameLoop := range ctx.Config.GetConnections(fullNodename) {
			nodeID := execute.ShowNodeID(ctx.Config.GetBinary(fullNodenameLoop), ctx.Config.GetHome(fullNodenameLoop))
			peers = append(peers, fmt.Sprintf("%s@127.0.0.1:%d\n", nodeID, ctx.Config.GetP2PPort(fullNodenameLoop)))
			peerIDs = append(peerIDs, nodeID)
		}
		utils.SetConfigEntry(configToml, "p2p.persistent_peers", strings.Join(peers, ","))
		utils.SetConfigEntry(configToml, "p2p.unconditional_peer_ids", strings.Join(peerIDs, ","))
		utils.SetConfigEntry(configToml, "p2p.external_address", fmt.Sprintf("127.0.0.1:%d", ctx.Config.GetP2PPort(fullNodename)))
	}
}
