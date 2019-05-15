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

	eg.POST("/users/:id", func(c *hp.Context) {
		c.JSON(c.Req.Param["id"], nil)
	})

	eg_1 := eg.Group("/test/", nil)
	{
		eg_1.POST("/Add/", userCtr.Add)
		eg_1.POST("/Update/", userCtr.Update)
		eg_1.GET("/Query/", userCtr.Query)
		eg_1.DELETE("/Delete/:id", userCtr.Delete)
	}
}

type Test struct {
	Name string
}
