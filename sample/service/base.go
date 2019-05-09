package service

import "github.com/jinzhu/gorm"
import "go-open/library/database/orm"

type BaseService struct {
	_DB *gorm.DB
}

func (base *BaseService) DB() *gorm.DB {

	if base._DB == nil {
		base._DB = orm.DB
	}
	return base._DB
}
