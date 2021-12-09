package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
	"math/rand"
	"net/http"
)

type Config struct {
	Siblings []string
	IsLeaf   bool
}

var config Config

func Init() {
	rand.Seed(17)
	envconfig.MustProcess("sibserver", &config)
	config.IsLeaf = len(config.Siblings) == 0

	log.Println("config", config.Siblings, config.IsLeaf)
}

func main() {
	Init()

	r := gin.Default()

	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	r.GET("/status", func(ctx *gin.Context) {
		rng := rand.Int31n(2) == 0
		if rng {
			ctx.Status(http.StatusOK)
		} else {
			ctx.Status(http.StatusTeapot)
		}
	})

	r.GET("/data", func(ctx *gin.Context) {
		if config.IsLeaf {
			ctx.Status(http.StatusOK)
			ctx.Next()
		}
		for _, sib := range config.Siblings {
			res, err := http.Get(sib)
			if err != nil {
				log.Println("err", err)
			}
			log.Println("sibling status", sib, res.Status)
		}
	})

	r.Run(":3440")
}
