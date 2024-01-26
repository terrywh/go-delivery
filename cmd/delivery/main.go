package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"

	apiv1 "github.com/terrywh/go-delivery/app/v1/api"
	dotv1 "github.com/terrywh/go-delivery/app/v1/dot"
)

var appHost host.Host
var appDHT *dht.IpfsDHT



func main() {
	var err error
	ctx := context.Background()
	log.Println("bootstrap peers: ", dht.DefaultBootstrapPeers)
	
	if appHost, err = libp2p.New(); err != nil {
		log.Fatal("failed to create host: ", err)
		return
	}
	log.Println("host: ", appHost.ID(), appHost.Addrs())
	appHost.SetStreamHandler(apiv1.Protocol, apiv1.StreamHandler)
	appHost.SetStreamHandler(dotv1.Protocol, dotv1.StreamHandler)

	if appDHT, err = dht.New(ctx, appHost); err != nil {
		log.Fatal("failed to create dht: ", err)
		return
	}
	appDHT.Bootstrap(ctx)
	

	var wg sync.WaitGroup
	for _, addr := range dht.DefaultBootstrapPeers {
		ai, _ := peer.AddrInfoFromP2pAddr(addr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := appHost.Connect(ctx, *ai); err != nil {
				log.Println("failed to connect to bootstrap peer: ", ai.ID, ai.Addrs, ", due to: ", err)
			} else {
				log.Println("connected to boostrap peer: ", ai.ID, ai.Addrs)
			}
		} ()
	}
	wg.Wait()

	discovery := routing.NewRoutingDiscovery(appDHT)
	util.Advertise(ctx, discovery, "github.com/terrywh/go-delivery")

	for {
		log.Println(".")
		peers, err := discovery.FindPeers(ctx, "github.com/terrywh/go-delivery")
		if err != nil {
			log.Fatal("failed to find peers: ", err)
			return
		}

		for peer := range peers {
			if peer.ID == appHost.ID() {
				continue
			}

			log.Println("found peer: ", peer.ID, peer.Addrs)
			// appHost.NewStream(ctx, peer.ID, dotv1.Protocol)
		}
		time.Sleep(15 * time.Second)
	}
}