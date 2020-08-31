package model

import (
	"github.com/xiaolongfan119/go-open/library/database/orm"
)

type Address struct {
	orm.Model
	Address string
}
