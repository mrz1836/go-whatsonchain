// Package main demonstrates API key usage with the WhatsOnChain client.
package main

import (
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create options (add your api key)
	opts := whatsonchain.ClientDefaultOptions()
	opts.APIKey = "your-secret-key"

	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, opts, nil)
	log.Println("client loaded", client.UserAgent())
}
