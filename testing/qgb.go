package testing

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/celestiaorg/orchestrator-relayer/helpers"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	celestiatestnode "github.com/celestiaorg/celestia-app/testutil/testnode"
	"github.com/celestiaorg/orchestrator-relayer/orchestrator"
	"github.com/celestiaorg/orchestrator-relayer/p2p"
	"github.com/celestiaorg/orchestrator-relayer/relayer"
	"github.com/celestiaorg/orchestrator-relayer/rpc"

	"github.com/celestiaorg/orchestrator-relayer/evm"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func NewRelayer(
	t *testing.T,
	node *TestNode,
) *relayer.Relayer {
	logger := tmlog.NewNopLogger()
	node.CelestiaNetwork.GRPCClient.Close()
	appQuerier := rpc.NewAppQuerier(logger, node.CelestiaNetwork.GRPCAddr, encoding.MakeConfig(app.ModuleEncodingRegisters...))
	require.NoError(t, appQuerier.Start())
	t.Cleanup(func() {
		_ = appQuerier.Stop()
	})
	tmQuerier := rpc.NewTmQuerier(node.CelestiaNetwork.RPCAddr, logger)
	tmQuerier.WithClientConn(node.CelestiaNetwork.Client)
	p2pQuerier := p2p.NewQuerier(node.DHTNetwork.DHTs[0], logger)
	evmClient := NewEVMClient(node.EVMChain.Key)
	r := relayer.NewRelayer(tmQuerier, appQuerier, p2pQuerier, evmClient, logger)
	return r
}

func NewEVMClient(key *ecdsa.PrivateKey) *evm.Client {
	logger := tmlog.NewNopLogger()
	// specifying an empty RPC endpoint as we will not be testing the methods that require it.
	// the simulated backend doesn't provide an RPC endpoint.
	return evm.NewClient(logger, nil, key, "", 100000000)
}

func NewOrchestrator(
	t *testing.T,
	node *TestNode,
) *orchestrator.Orchestrator {
	logger := tmlog.NewNopLogger()
	appQuerier := rpc.NewAppQuerier(logger, node.CelestiaNetwork.GRPCAddr, encoding.MakeConfig(app.ModuleEncodingRegisters...))
	require.NoError(t, appQuerier.Start())
	t.Cleanup(func() {
		_ = appQuerier.Stop()
	})
	tmQuerier := rpc.NewTmQuerier(node.CelestiaNetwork.RPCAddr, logger)
	tmQuerier.WithClientConn(node.CelestiaNetwork.Client)
	p2pQuerier := p2p.NewQuerier(node.DHTNetwork.DHTs[0], logger)
	broadcaster := orchestrator.NewBroadcaster(node.DHTNetwork.DHTs[0])
	retrier := helpers.NewRetrier(logger, 3, 500*time.Millisecond)
	orch, err := orchestrator.New(logger, appQuerier, tmQuerier, p2pQuerier, broadcaster, retrier, *celestiatestnode.NodeEVMPrivateKey)
	require.NoError(t, err)
	return orch
}
