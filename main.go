package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

// Config : captures app config from env
type Config struct {
	CoinMarketCap struct {
		Key string `envconfig:"CMC_API_KEY"`
	}
	AzureBlob struct {
	}
	CosmosDB struct {
	}
}

func main() {
	var (
		// cli flags
		local     bool
		useblob   bool
		usecosmos bool
		verbose   bool
		limit     string
		start     string

		// app config
		cfg Config
	)

	_ = envconfig.Process("", &cfg)

	flag.BoolVar(&local, "l", true, "Read items from local disk.")
	flag.BoolVar(&useblob, "use-blob", false, "Write items to Azure Blob.")
	flag.BoolVar(&usecosmos, "use-cosmos", false, "Write items to CosmosDB.")
	flag.BoolVar(&verbose, "v", false, "Verbose logging.")
	flag.StringVar(&limit, "limit", "100", "Specify the number of results to return. Use this parameter and the \"start\" parameter to determine your own pagination size.")
	flag.StringVar(&start, "start", "1", "Offset the start (1-based index) of the paginated list of items to return.")
	flag.Parse()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
