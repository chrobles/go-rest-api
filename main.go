package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/chrobles/go-rest-api/cmclient"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	var (
		// coinmarketcap API client
		cmClient cmclient.Client
		apiKey   string

		// cli flags
		local     bool
		useblob   bool
		usecosmos bool
		verbose   bool
		limit     string
		start     string

		// response data
		res cmclient.RangeData
	)

	flag.BoolVar(&local, "l", true, "Read items from local disk.")
	flag.BoolVar(&useblob, "use-blob", false, "Write items to Azure Blob.")
	flag.BoolVar(&usecosmos, "use-cosmos", false, "Write items to CosmosDB.")
	flag.BoolVar(&verbose, "v", false, "Verbose logging.")
	flag.StringVar(&limit, "limit", "100", "Specify the number of results to return. Use this parameter and the \"start\" parameter to determine your own pagination size.")
	flag.StringVar(&start, "start", "1", "Offset the start (1-based index) of the paginated list of items to return.")
	flag.Parse()

	if !local {
		// cm client configuration
		apiKey = os.Getenv("CM_API_KEY")
		if apiKey == "" {
			log.Print("CM_API_KEY not found")
			os.Exit(1)
		}
		cmClient.Key = apiKey
		cmClient.Address = "https://sandbox-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
		req := cmClient.NewRangeRequest(start, limit)
		res = cmClient.Get(req)
	} else {
		data, _ := ioutil.ReadFile("data.json")
		json.Unmarshal(data, &res)
	}

	if dump {
		spew.Dump(res)
	}

	if upsert {
		spew.Dump(upsert)
	}
}
