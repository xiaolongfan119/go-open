package main

import (
	ecode "go-open/library/ecode"
	hp "go-open/library/net/http/hypnus"
	"go-open/sample/conf"
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

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	// })

	// http.ListenAndServe("")
}
