// Package main demonstrates API key usage with the WhatsOnChain client.
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client with API key
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
		whatsonchain.WithAPIKey("your-secret-key"),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client loaded", client.UserAgent())
}
