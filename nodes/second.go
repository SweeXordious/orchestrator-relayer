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

	decKey, err := hex.DecodeString("1BFC789DBD7B3CA13B4CF47898088CBB5CE467668DA63740ADF62B06F474452C6E12BD8B0C964D17438B8FEE1AC019D5290E2D4BE5BEED0113E13926581FFCB4")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	privKey, err := crypto.UnmarshalEd25519PrivateKey(decKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	m, err := multiaddr.NewMultiaddr("/ip4/192.168.0.111/tcp/37752")
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
	fmt.Println("created node2:")
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
	err = n1.Bstore.Put(ctx, blk)
	if err != nil {
		return
	}
	fmt.Printf("node putting block with CID: %s\n", blk.Cid().String())

	time.Sleep(5 * time.Second)

	err = n1.BitswapExchange.NotifyNewBlocks(ctx, blk)
	if err != nil {
		return
	}
	fmt.Println("notified for block")

	thirdId, err := peer.Decode("12D3KooWHMJPb8C69aX4VczGAdKmc8m9ZaV4qc15RNgqrXc4vSna")
	for {
		stat, err := n1.BitswapExchange.Stat()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println()
		fmt.Printf("%+v", stat)
		fmt.Println()
		fmt.Println(n1.BitswapExchange.WantlistForPeer(thirdId))
		time.Sleep(10 * time.Second)
	}
}
