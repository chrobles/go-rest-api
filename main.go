package main

import (
	"log"
	"strconv"

	"github.com/chrobles/go-rest-api/cmcapiclient"
	"github.com/chrobles/go-rest-api/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	cfg types.Config
)

func init() {
	// defaults
	cfg.CoinMarketCap.BaseURL = "https://sandbox-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
	cfg.CoinMarketCap.UseLocal = true

	// load vars from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	// parse config to struct
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	var (
		cmcclient cmcapiclient.Client
		err       error
	)

	err = cmcclient.Configure(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()
	r.GET("/getcmc/:limit", func(c *gin.Context) {
		var (
			limit   int
			start   int
			blob    bool
			cosmos  bool
			mktdata *types.MarketListings
		)

		limit, _ = strconv.Atoi(c.Param("limit"))
		start, _ = strconv.Atoi(c.DefaultQuery("start", "1"))
		blob, _ = strconv.ParseBool(c.DefaultQuery("useblob", "false"))
		cosmos, _ = strconv.ParseBool(c.DefaultQuery("usecosmos", "false"))

		mktdata, err = cmcclient.GetMarketListings(start, limit)
		if err != nil {
			log.Print(err)
			c.JSON(400, gin.H{
				"error": err,
			})
		} else {
			c.JSON(200, mktdata)
		}

		if blob == true {
		}
		if cosmos == true {
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
