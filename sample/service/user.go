package service

import (
	"go-open/sample/model"
)

type UserService struct {
	BaseService
}

func (service *UserService) Add() (data interface{}, err error) {
	user := &model.User{Age: 10, Name: "hai"}
	service.DB().Create(user)
	panic("errrrrr")

	return
}
