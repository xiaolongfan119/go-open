package http

import (
	"net/http"

	"github.com/xiaolongfan119/go-open/v2/sample/controller"

	hp "github.com/xiaolongfan119/go-open/v2/library/net/http/hypnus"
)

var (
	userCtr = &controller.UserController{}

	fileCtr = &controller.FileController{}
)

func Init(conf *hp.ServerConf) {
	engine := hp.DefaultServer(conf)
	route(engine)
	engine.Start()
}

func route(engine *hp.Engine) {

	engine.ServeFiles("/static/*filepath", http.Dir("./../assets"))

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

	eg_2 := eg.Group("/file/", nil)
	{
		eg_2.POST("/upload", fileCtr.Upload)
	}
}

type Test struct {
	Name string
}
