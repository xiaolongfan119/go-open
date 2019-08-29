package model

import (
	"github.com/ihornet/go-open/library/database/orm"
)

type Address struct {
	orm.Model
	Address string
}
