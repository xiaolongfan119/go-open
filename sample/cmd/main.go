package main

import (
	hp "go-open/library/net/http/hypnus"
	"go-open/sample/conf"

	jwt "go-open/library/net/http/hypnus/middleware/token"
)

func main() {
	conf.Init()
	token := jwt.Token{Conf: conf.Conf.TokenConfig}
	engine := hp.DefaultServer(conf.Conf.HttpServer)
	eg := engine.Group("/", token.Verify)
	eg.POST("/test", func(c *hp.Context) {
		// c.Writer.Write([]byte(c.Req.Body["name"].(string)))

		//c.JSON(c.Req.Body["name"], nil)

		c.JSON(token.GenToken(nil), nil)

		//	c.JSON(nil, ecode.UnknownCode)
	})
	engine.Start()
}
