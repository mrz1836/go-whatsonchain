// Package main demonstrates retrieving UTXOs for an address using the whatsonchain client.
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get UTXOs for an address
	history, err := client.AddressUnspentTransactions(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	if err != nil {
		log.Printf("error getting utxos: %s", err.Error())
	} else if len(history) == 0 {
		log.Println("no utxos found")
	} else {
		for index, utxo := range history {
			log.Printf("(%d) %s | Sats: %d \n", index+1, utxo.TxHash, utxo.Value)
		}
	}
}
