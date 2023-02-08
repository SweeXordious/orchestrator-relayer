package testing

import (
	"crypto/ecdsa"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/orchestrator-relayer/evm"
	"github.com/celestiaorg/orchestrator-relayer/orchestrator"
	"github.com/celestiaorg/orchestrator-relayer/p2p"
	"github.com/celestiaorg/orchestrator-relayer/relayer"
	"github.com/celestiaorg/orchestrator-relayer/rpc"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func NewOrchestrator(
	node *TestNode,
) *orchestrator.Orchestrator {
	logger := tmlog.NewNopLogger()
	appQuerier := rpc.NewAppQuerier(logger, node.CelestiaNetwork.GRPCClient, encoding.MakeConfig(app.ModuleEncodingRegisters...))
	tmQuerier := rpc.NewTmQuerier(node.CelestiaNetwork.Client, logger)
	p2pQuerier := p2p.NewQuerier(node.DHTNetwork.DHTs[0], logger)
	broadcaster := orchestrator.NewBroadcaster(node.DHTNetwork.DHTs[0])
	retrier := orchestrator.NewRetrier(logger, 3)
	orch, err := orchestrator.New(logger, appQuerier, tmQuerier, p2pQuerier, broadcaster, retrier, *node.EVMChain.Key)
	if err != nil {
		panic(err)
	}
	return orch
}

func NewRelayer(
	node *TestNode,
) *relayer.Relayer {
	logger := tmlog.NewNopLogger()
	appQuerier := rpc.NewAppQuerier(logger, node.CelestiaNetwork.GRPCClient, encoding.MakeConfig(app.ModuleEncodingRegisters...))
	tmQuerier := rpc.NewTmQuerier(node.CelestiaNetwork.Client, logger)
	p2pQuerier := p2p.NewQuerier(node.DHTNetwork.DHTs[0], logger)
	evmClient := NewEVMClient(node.EVMChain.Key)
	r, err := relayer.NewRelayer(tmQuerier, appQuerier, p2pQuerier, evmClient, logger)
	if err != nil {
		panic(err)
	}
	return r
}

func NewEVMClient(key *ecdsa.PrivateKey) *evm.Client {
	logger := tmlog.NewNopLogger()
	// specifying an empty RPC as we will not be testing the methods that require it.
	// the simulated backend doesn't provide an RPC endpoint.
	return evm.NewClient(logger, nil, key, "", 100000000)
}
