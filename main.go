package main

import (
	"log"
	"os"

	"github.com/chrobles/go-rest-api/cmclient"
)

var (
	cmClient cmclient.Client
)

func main() {
	apiKey := os.Getenv("CM_API_KEY")

	if apiKey == "" {
		log.Print("CM_API_KEY not found")
		os.Exit(1)
	}

	cmClient.Key = apiKey
	cmClient.Address = "https://sandbox-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"

	req := cmClient.NewRangeRequest(os.Args[1], os.Args[2])
	cmClient.Get(req)
}
