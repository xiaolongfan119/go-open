package controller

import (
	hp "go-open/library/net/http/hypnus"
	"go-open/sample/service"
)

var userSrv = &service.UserService{}

type UserController struct{}

func (ctr *UserController) Add(ctx *hp.Context) {
	ctx.JSON(userSrv.Add())
}
