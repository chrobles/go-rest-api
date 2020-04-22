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
	cfg.CosmosDB.UseCosmos = false

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

	if cfg.CosmosDB.UseCosmos {
		err = cdbclient.Configure(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
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

		if cosmos && cfg.CosmosDB.UseCosmos == false {
			c.JSON(400, gin.H{"error": "CDB_USE_COSMOS = false"})
		}

		mktdata, err = cmcclient.GetMarketListings(start, limit)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		} else if !cosmos {
			c.JSON(200, mktdata)
		}

		if cosmos && cfg.CosmosDB.UseCosmos == true {
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
