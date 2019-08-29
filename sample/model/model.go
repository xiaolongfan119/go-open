package model

import (
	"github.com/ihornet/go-open/library/database/orm"

	"github.com/ihornet/go-open/sample/conf"
)

func Init(conf *conf.Config) {

	models := []interface{}{&User{}, &Address{}, &Email{}, &Profile{}}

	db := orm.NewMySQL(conf.DBConfig)
	db.AutoMigrate(models...)
}
