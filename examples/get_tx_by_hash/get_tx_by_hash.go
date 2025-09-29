// Package main demonstrates transaction lookup by hash functionality.
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	// Get the transaction information
	info, _ := client.GetTxByHash(context.Background(), "908c26f8227fa99f1b26f99a19648653a1382fb3b37b03870e9c138894d29b3b")
	log.Printf("block hash: %s", info.BlockHash)
}
