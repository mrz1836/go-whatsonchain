package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	// Get the balance for multiple addresses
	balances, _ := client.BulkUnspentTransactions(
		context.Background(),
		&whatsonchain.AddressList{
			Addresses: []string{
				"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
				"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob",
			},
		},
	)

	for _, record := range balances {
		log.Printf(
			"address: %s utxos: %d \n",
			record.Address,
			len(record.Utxos),
		)
	}
}
