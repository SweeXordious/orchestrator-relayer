package rpc

import (
	"context"

	"github.com/tendermint/tendermint/libs/bytes"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// TmQuerierI queries tendermint for commitments and events.
type TmQuerierI interface {
	// QueryCommitment queries tendermint for a commitment for the set of blocks
	// defined by `beginBlock` and `endBlock`.
	QueryCommitment(ctx context.Context, beginBlock uint64, endBlock uint64) (bytes.HexBytes, error)

	// SubscribeEvents subscribe to the events defined by the provided query.
	SubscribeEvents(ctx context.Context, subscriptionName string, query string) (<-chan coretypes.ResultEvent, error)
}

var _ TmQuerierI = &TmQuerier{}

type TmQuerier struct {
	logger        tmlog.Logger
	TendermintRPC client.Client
}

func NewTmQuerier(
	tendermintRPC client.Client,
	logger tmlog.Logger,
) *TmQuerier {
	return &TmQuerier{
		logger:        logger,
		TendermintRPC: tendermintRPC,
	}
}

func (tq TmQuerier) QueryCommitment(ctx context.Context, beginBlock uint64, endBlock uint64) (bytes.HexBytes, error) {
	dcResp, err := tq.TendermintRPC.DataCommitment(ctx, beginBlock, endBlock)
	if err != nil {
		return nil, err
	}
	return dcResp.DataCommitment, nil
}

func (tq TmQuerier) SubscribeEvents(ctx context.Context, subscriptionName string, query string) (<-chan coretypes.ResultEvent, error) {
	// This doesn't seem to complain when the node is down
	results, err := tq.TendermintRPC.Subscribe(
		ctx,
		subscriptionName,
		query,
	)
	if err != nil {
		return nil, err
	}
	return results, err
}
