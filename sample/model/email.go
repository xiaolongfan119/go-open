package model

import (
	"github.com/ihornet/go-open/library/database/orm"
)

type Email struct {
	orm.Model
	Email  string `json:"email"`
	UserID uint   `json:"userID"`
}
