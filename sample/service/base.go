package service

import (
	"github.com/ihornet/go-open/v2/library/database/orm"
	"github.com/jinzhu/gorm"
)

type BaseService struct {
	_DB *gorm.DB
}

func (base *BaseService) DB() *gorm.DB {

	if base._DB == nil {
		base._DB = orm.DB
	}
	return base._DB
}
