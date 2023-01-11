package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/celestiaorg/orchestrator-relayer/p2p"
	blocks "github.com/ipfs/go-block-format"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"time"
)

func main() {
	ctx := context.Background()

	decKey, err := hex.DecodeString("6A0AEC628DAAA68A49A047D1DFB38D1A882F1DBA69E52CB6B595AB12BB59E1D96FEE8A894E21F3EC651FBE7CB93DCB94B94A30D79E4D413CBC9BFE3B315D912F")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	privKey, err := crypto.UnmarshalEd25519PrivateKey(decKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	m, err := multiaddr.NewMultiaddr("/ip4/192.168.0.111/tcp/37753")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	h, err := libp2p.New(libp2p.ListenAddrs(m), libp2p.Identity(privKey), libp2p.EnableNATService())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	n1, err := p2p.CreateNode(ctx, h)
	if err != nil {
		return
	}
	fmt.Println("created node3:")
	fmt.Println(n1.Host.ID())
	fmt.Println(n1.Host.Addrs())

	Id, err := peer.Decode("12D3KooWBxJm9dVuxL3Yiwvx2BDgRSJgpw7kbFtspa8ns4EWzpj1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	multiAddr, err := multiaddr.NewMultiaddr("/ip4/192.168.0.111/tcp/37751")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("connecting to %s\n", Id)

	err = n1.Router.Host().Connect(ctx, peer.AddrInfo{
		ID:    Id,
		Addrs: []multiaddr.Multiaddr{multiAddr},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("connected node to %s\n", Id)

	blk := blocks.NewBlock([]byte("hello"))
	fmt.Printf("block CID: %s\n", blk.Cid().String())

	go func() {
		for i := 0; i < 5; i++ {
			stat, err := n1.BitswapExchange.Stat()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println()
			fmt.Printf("%+v", stat)
			fmt.Println()
			time.Sleep(10 * time.Second)
		}
	}()

	block, err := n1.BitswapExchange.GetBlock(ctx, blk.Cid())
	if err != nil {
		return
	}

	fmt.Println(block.RawData())
	time.Sleep(60 * time.Second)
}
