package model

import (
	"go-open/library/database/orm"
)

type Address struct {
	orm.Model
	Address string
}
