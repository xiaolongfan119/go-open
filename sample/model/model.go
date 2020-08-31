package model

import (
	"github.com/xiaolongfan119/go-open/library/database/orm"

	"github.com/xiaolongfan119/go-open/sample/conf"
)

func Init(conf *conf.Config) {

	models := []interface{}{&User{}, &Address{}, &Email{}, &Profile{}}

	db := orm.NewMySQL(conf.DBConfig)
	db.AutoMigrate(models...)
}
