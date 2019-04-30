package main

import (
	ecode "github.com/ihornet/go-commom/library/ecode"
	hp "github.com/ihornet/go-commom/library/net/http/hypnus"
	"github.com/ihornet/go-commom/sample/conf"
)

func main() {
	conf.Init()
	engine := hp.DefaultServer(conf.Conf.HttpServer)
	eg := engine.Group("/", nil)
	eg.POST("/test", func(c *hp.Context) {
		// c.Writer.Write([]byte(c.Req.Body["name"].(string)))
		//c.JSON(c.Req.Body["name"].(string), nil)
		c.JSON(nil, ecode.UnknownCode)
	})
	engine.Start()
}
