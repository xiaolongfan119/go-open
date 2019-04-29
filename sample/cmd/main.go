package main

import (
	"fmt"

	"github.com/ihornet/go-commom/library/net/http/hypnus"
	"github.com/ihornet/go-commom/sample/conf"
)

func main() {
	conf.Init()
	engine := hypnus.DefaultServer(conf.Conf.HttpServer)
	eg := engine.Group("/", nil)
	eg.GET("/test", func(c *hypnus.Context) {
		fmt.Println("##########")
	})
}
