package main

import (
	log "go-open/library/log"
	hp "go-open/library/net/http/hypnus"
	"go-open/sample/conf"
)

func main() {
	conf.Init()
	//token := jwt.Token{Conf: conf.Conf.TokenConfig}
	log.Init(conf.Conf.LogConfig)
	log.Error("#######")
	engine := hp.DefaultServer(conf.Conf.HttpServer)
	eg := engine.Group("/", nil)
	eg.POST("/test/:id", func(c *hp.Context) {
		// c.Writer.Write([]byte(c.Req.Body["name"].(string)))

		c.JSON(c.Req.Param["id"], nil)

		//c.JSON(token.GenToken(nil), nil)
		//	c.JSON(nil, ecode.UnknownCode)
	})

	eg.POST("/test/", func(c *hp.Context) {
		c.JSON("hahh", nil)
	})

	eg.POST("/test2/", func(c *hp.Context) {
		c.JSON("test2", nil)
	})

	eg.GET("/test2/", func(c *hp.Context) {
		c.JSON("get  test2", nil)
	})

	engine.Start()

}
