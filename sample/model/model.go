package model

import (
	"github.com/xiaolongfan119/go-open/v2/library/database/orm"

	"github.com/xiaolongfan119/go-open/v2/sample/conf"
)

func Init(conf *conf.Config) {

	models := []interface{}{&User{}, &Address{}, &Email{}, &Profile{}}

	db := orm.NewMySQL(conf.DBConfig)
	db.AutoMigrate(models...)
}
