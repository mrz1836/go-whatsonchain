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

	// Create a client with custom HTTP client.
	//
	// A custom HTTP client is used as is, so the SDK's retries, backoff,
	// Retry-After handling and request timeout do not apply to it. To customize
	// only the transport (a proxy, TLS settings, tracing) and keep retries, use
	// whatsonchain.WithRoundTripper(http.DefaultTransport) instead.
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
