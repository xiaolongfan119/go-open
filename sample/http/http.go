package http

import (
	hp "go-open/library/net/http/hypnus"
	"go-open/sample/controller"
)

var (
	userCtr = &controller.UserController{}
)

func Init(conf *hp.ServerConf) {
	engine := hp.DefaultServer(conf)
	route(engine)
	engine.Start()
}

func route(engine *hp.Engine) {
	eg := engine.Group("/", nil)
	eg.POST("/test/:id", func(c *hp.Context) {
		c.JSON(c.Req.Param["id"], nil)
	})

	eg.POST("/Add/", userCtr.Add)
	eg.POST("/Update/", userCtr.Update)
	eg.GET("/Query/", userCtr.Query)
	eg.DELETE("/Delete/:id", userCtr.Delete)

}

type Test struct {
	Name string
}
