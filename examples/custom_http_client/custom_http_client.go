package main

import (
	"log"
	"net/http"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// use your own custom http client
	customClient := http.DefaultClient

	// Create a client
	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, customClient)
	log.Println("client loaded", client.UserAgent())
}
