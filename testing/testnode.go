package testing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
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

	// TODO check if this is correct
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	privateKey, err := crypto.HexToECDSA(hex.EncodeToString(bytes))
	if err != nil {
		panic(err)
	}

	evmChain := NewEVMChain(privateKey)

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
