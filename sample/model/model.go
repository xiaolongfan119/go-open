package model

import (
	"go-open/library/database/orm"
	"go-open/sample/conf"
)

func Init(conf *conf.Config) {

	models := []interface{}{&User{}, &Address{}}

	db := orm.NewMySQL(conf.DBConfig)
	db.AutoMigrate(models...)
}
