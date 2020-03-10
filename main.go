package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chrobles/go-rest-api/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	cfg types.Config
)

func init() {
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
	r := gin.Default()
	r.GET("/getcmc/:limit", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.Param("limit"))
		start, _ := strconv.Atoi(c.DefaultQuery("start", "1"))
		blob, _ := strconv.ParseBool(c.DefaultQuery("useblob", "false"))
		cosmos, _ := strconv.ParseBool(c.DefaultQuery("usecosmos", "false"))

		fmt.Print(limit, blob, cosmos)

		c.JSON(200, gin.H{
			"limit":     limit,
			"start":     start,
			"useblob":   blob,
			"usecosmos": cosmos,
			"config":    cfg,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
