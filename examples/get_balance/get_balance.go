package main

import (
	"context"
	"fmt"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {

	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	// Get a balance for an address
	balance, _ := client.AddressBalance(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	fmt.Println("confirmed balance", balance.Confirmed)
}
