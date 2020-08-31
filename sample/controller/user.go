package controller

import (
	"github.com/xiaolongfan119/go-open/sample/model"
	"github.com/xiaolongfan119/go-open/sample/service"

	hp "github.com/xiaolongfan119/go-open/library/net/http/hypnus"
)

var userSrv = &service.UserService{}

type UserController struct{}

func (ctr *UserController) Add(ctx *hp.Context) {

	user := new(model.User)
	if err := ctx.Bind(user); err != nil {
		return
	}
	ctx.JSON(userSrv.Add(user))
}

func (ctr *UserController) Update(ctx *hp.Context) {
	ctx.JSON(userSrv.Update())
}

func (ctr *UserController) Query(ctx *hp.Context) {
	ctx.JSON(userSrv.Query())
}

func (ctr *UserController) Delete(ctx *hp.Context) {
	ctx.JSON(userSrv.Delete(ctx.Req.Param["id"]))
}
