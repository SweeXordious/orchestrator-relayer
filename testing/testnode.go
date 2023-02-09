package testing

import (
	"context"
	"testing"
	"time"

	celestiatestnode "github.com/celestiaorg/celestia-app/testutil/testnode"
)

// TestNode contains a DHTNetwork along with a test Celestia network.
type TestNode struct {
	Context         context.Context
	DHTNetwork      *DHTNetwork
	CelestiaNetwork *CelestiaNetwork
	EVMChain        *EVMChain
}

func NewTestNode(ctx context.Context, t *testing.T) *TestNode {
	celestiaNetwork := NewCelestiaNetwork(t, time.Millisecond)
	// minimum number of peers for a DHT is 2. If not, it will not be able to put values.
	dhtNetwork := NewDHTNetwork(ctx, 2)

	evmChain := NewEVMChain(celestiatestnode.NodeEVMPrivateKey)
	go evmChain.PeriodicCommit(ctx, DEFAULT_PERIODIC_COMMIT_DELAY)

	return &TestNode{
		Context:         ctx,
		DHTNetwork:      dhtNetwork,
		CelestiaNetwork: celestiaNetwork,
		EVMChain:        evmChain,
	}
}

func (tn TestNode) Close() {
	tn.DHTNetwork.Stop()
	tn.CelestiaNetwork.Stop()
	tn.EVMChain.Close()
}
