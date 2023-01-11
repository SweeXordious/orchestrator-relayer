package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/celestiaorg/orchestrator-relayer/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"time"
)

func main() {
	ctx := context.Background()

	multiAddr, err := multiaddr.NewMultiaddr("/ip4/192.168.0.111/tcp/37751")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	decKey, err := hex.DecodeString("398E0B0478862529D79F21C028317DD181C1E67AF44D099CDC78BD064A13DF671FC02B117A792377BFE4E931045FB47BAB9B236ED3AC43A254A76E9CE9CAD8B8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	privKey, err := crypto.UnmarshalEd25519PrivateKey(decKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	h, err := libp2p.New(libp2p.ListenAddrs(multiAddr), libp2p.Identity(privKey), libp2p.EnableNATService())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	n1, err := p2p.CreateNode(ctx, h)
	if err != nil {
		return
	}
	fmt.Println("created nodes:")
	fmt.Println(n1.Host.ID())
	fmt.Println(n1.Host.Addrs())

	thirdId, err := peer.Decode("12D3KooWHMJPb8C69aX4VczGAdKmc8m9ZaV4qc15RNgqrXc4vSna")
	if err != nil {
		return
	}

	for i := 0; i < 15; i++ {
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

	time.Sleep(60 * time.Second)
}
