package main

import (
	"log"
	"time"

	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

func createNodeConnection() api.BitsharesAPI {
	connected := make(chan struct{})
	conn, rpc := api.NewBuilder().Node(cfg.nodeAddress).
		LoginHandler(func() { connected <- struct{}{} }).Build()
	conn.Connect()

	select {
	case <-connected:
		return rpc
	case <-time.After(5 * time.Second):
		log.Fatalf("Connection to %s timed out", cfg.nodeAddress)
	}

	return nil
}

func getAsset(api api.DatabaseAPI, symbol string) objects.Asset {
	assets, err := api.ListAssets(listOptions.asset, 1)
	if err != nil {
		log.Fatalf("Failed to get asset: %v", err)
	}
	if len(assets) < 1 || assets[0].Symbol != listOptions.asset {
		log.Fatalf("Asset %s not found", listOptions.asset)
	}

	return assets[0]
}
