package service

import (
	"github.com/xiaolongfan119/go-open/v2/sample/model"
)

type UserService struct {
	BaseService
}

func (service *UserService) Add(user *model.User) (data interface{}, err error) {
	profile := model.Profile{Name: "ppsssp"}
	email := model.Email{Email: "email###sads##"}
	user.Profile = profile
	user.Email = email
	ret := service.DB().Create(&user)

	return ret.Value, nil
}

func (service *UserService) Update() (data interface{}, err error) {
	user := model.User{Age: 11}
	user.ID = 1
	//	user.Profile = model.Profile{Name: "#######"}
	//	service.DB().Model(&user).Select("name").Updates(map[string]interface{}{"name": "hahhah"})
	service.DB().Save(&user).Update("name", "hello")
	return user, nil
}

func (service *UserService) Query() (data interface{}, err error) {
	var user model.User
	user.ID = 10
	ret := service.DB().Preload("Email").Preload("Profile").Find(&user)
	err = ret.Error
	return user, err
}

func (service *UserService) Delete(id string) (data interface{}, err error) {
	service.DB().Delete(model.User{}, "id=?", id)
	return
}
