package model

import (
	"github.com/xiaolongfan119/go-open/v2/library/database/orm"
)

type Address struct {
	orm.Model
	Address string
}
