package orchestrator_test

import (
	"context"
	"testing"

	"github.com/celestiaorg/orchestrator-relayer/orchestrator"
	"github.com/celestiaorg/orchestrator-relayer/p2p"
	qgbtesting "github.com/celestiaorg/orchestrator-relayer/testing"
	"github.com/celestiaorg/orchestrator-relayer/types"
	"github.com/stretchr/testify/assert"
)

func TestBroadcastDataCommitmentConfirm(t *testing.T) {
	network := qgbtesting.NewDHTNetwork(context.Background(), 4)
	defer network.Stop()

	// create a test DataCommitmentConfirm
	expectedConfirm := types.DataCommitmentConfirm{
		EthAddress: "celes1qktu8009djs6uym9uwj84ead24exkezsaqrmn5",
		Commitment: "test commitment",
		Signature:  "test signature",
	}

	// generate a test key for the DataCommitmentConfirm
	testKey := p2p.GetDataCommitmentConfirmKey(10, "celes1qktu8009djs6uym9uwj84ead24exkezsaqrmn5")

	// Broadcast the confirm
	broadcaster := orchestrator.NewBroadcaster(network.DHTs[1])
	err := broadcaster.BroadcastDataCommitmentConfirm(context.Background(), 10, expectedConfirm)
	assert.NoError(t, err)

	// try to get the confirm from another peer
	actualConfirm, err := network.DHTs[3].GetDataCommitmentConfirm(context.Background(), testKey)
	assert.NoError(t, err)
	assert.NotNil(t, actualConfirm)

	assert.Equal(t, expectedConfirm, actualConfirm)
}

func TestBroadcastValsetConfirm(t *testing.T) {
	network := qgbtesting.NewDHTNetwork(context.Background(), 4)
	defer network.Stop()

	// create a test DataCommitmentConfirm
	expectedConfirm := types.ValsetConfirm{
		EthAddress: "celes1qktu8009djs6uym9uwj84ead24exkezsaqrmn5",
		Signature:  "test signature",
	}

	// generate a test key for the ValsetConfirm
	testKey := p2p.GetValsetConfirmKey(10, "celes1qktu8009djs6uym9uwj84ead24exkezsaqrmn5")

	// Broadcast the confirm
	broadcaster := orchestrator.NewBroadcaster(network.DHTs[1])
	err := broadcaster.BroadcastValsetConfirm(context.Background(), 10, expectedConfirm)
	assert.NoError(t, err)

	// try to get the confirm from another peer
	actualConfirm, err := network.DHTs[3].GetValsetConfirm(context.Background(), testKey)
	assert.NoError(t, err)
	assert.NotNil(t, actualConfirm)

	assert.Equal(t, expectedConfirm, actualConfirm)
}
