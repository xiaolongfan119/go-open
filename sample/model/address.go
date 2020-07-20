package model

import (
	"github.com/ihornet/go-open/v2/library/database/orm"
)

type Address struct {
	orm.Model
	Address string
}
