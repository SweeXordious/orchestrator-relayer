package p2p

import (
	"context"
	"github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

func CreateNode(ctx context.Context, h host.Host) (*Node, error) {
	dstore := dssync.MutexWrap(ds.NewMapDatastore())
	bstore := blockstore.NewBlockstore(dstore)

	router, err := dht.New(
		ctx,
		h,
		dht.Datastore(dstore),
	)
	if err != nil {
		return nil, err
	}

	network := bsnet.NewFromIpfsHost(h, router)
	exchange := bitswap.New(
		ctx,
		network,
		bstore,
		//bitswap.ProvideEnabled(false),
		//bitswap.SetSendDontHaves(false),
	)

	return &Node{
		Host:            h,
		Router:          router,
		Bstore:          bstore,
		Network:         network,
		BitswapExchange: exchange,
	}, nil
}

type Node struct {
	Host            host.Host
	Router          *dht.IpfsDHT
	Bstore          blockstore.Blockstore
	Network         bsnet.BitSwapNetwork
	BitswapExchange *bitswap.Bitswap
}
