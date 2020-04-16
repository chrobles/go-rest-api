package main

import (
	"log"
	"os"
	"strconv"

	"github.com/chrobles/go-rest-api/cmcapiclient"
	"github.com/chrobles/go-rest-api/cosmosdbclient"
	"github.com/chrobles/go-rest-api/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	cfg types.Config
)

func init() {
	var (
		pwd string
		err error
	)

	pwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}

	// project root dir
	cfg.Root = pwd

	// defaults
	cfg.CoinMarketCap.BaseURL = "https://sandbox-api.coinmarketcap.com"
	cfg.CoinMarketCap.UseLocal = true

	// load vars from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	// parse config to struct
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	spew.Dump(cfg)
}

func main() {
	var (
		cdbclient cosmosdbclient.Client
		cmcclient cmcapiclient.Client
		err       error
	)

	err = cmcclient.Configure(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cdbclient.Configure(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()
	r.GET("/mkt/:limit", func(c *gin.Context) {
		var (
			err     error
			limit   int
			start   int
			cosmos  bool
			mktdata *types.MarketListings
		)

		limit, _ = strconv.Atoi(c.Param("limit"))
		start, _ = strconv.Atoi(c.DefaultQuery("start", "1"))
		cosmos, _ = strconv.ParseBool(c.DefaultQuery("cosmos", "false"))

		mktdata, err = cmcclient.GetMarketListings(start, limit)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else if !cosmos {
			c.JSON(200, mktdata)
		}

		if cosmos {
			var (
				id  interface{}
				err error
			)
			id, err = cdbclient.Index(mktdata)
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(200, id)
			}
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
