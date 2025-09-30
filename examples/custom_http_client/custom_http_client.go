// Package main demonstrates using a custom HTTP client with the whatsonchain library.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// use your own custom http client
	customClient := http.DefaultClient

	// Create a client with custom HTTP client
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
		whatsonchain.WithHTTPClient(customClient),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client loaded", client.UserAgent())
}
