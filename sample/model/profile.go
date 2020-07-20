package model

import (
	"github.com/ihornet/go-open/v2/library/database/orm"
)

type Profile struct {
	orm.Model
	Name   string `json:"name"`
	UserID uint   `json:"userID"`
}
