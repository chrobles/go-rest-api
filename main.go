package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

// Config : captures app config from env
type Config struct {
	App struct {
		Local bool `envconfig:"USE_LOCAL"`
	}
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

		// app config
		cfg Config
	)

	_ = envconfig.Process("", &cfg)

	flag.BoolVar(&local, "l", true, "Read items from local disk.")
	flag.BoolVar(&useblob, "use-blob", false, "Write items to Azure Blob.")
	flag.BoolVar(&usecosmos, "use-cosmos", false, "Write items to CosmosDB.")
	flag.BoolVar(&verbose, "v", false, "Verbose logging.")
	flag.Parse()

	r := gin.Default()
	r.GET("/getcmc/:limit", func(c *gin.Context) {
		limit, _ := strconv.ParseInt(c.Param("limit"))
		start, _ := strconv.ParseInt(c.DefaultQuery("start", "1"))
		blob, _ := strconv.ParseBool(c.DefaultQuery("useblob", "false"))
		cosmos, _ := strconv.ParseBool(c.DefaultQuery("usecosmos", "false"))

		fmt.Print(limit, blob, cosmos)

		c.JSON(200, gin.H{
			"limit":     limit,
			"useblob":   blob,
			"usecosmos": cosmos,
			"config":    cfg,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
