package model

import (
	"go-open/library/database/orm"
)

type Profile struct {
	orm.Model
	Name   string `json:"name"`
	UserID uint   `json:"userID"`
}
