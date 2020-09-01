package model

import (
	"github.com/xiaolongfan119/go-open/v2/library/database/orm"
)

type Profile struct {
	orm.Model
	Name   string `json:"name"`
	UserID uint   `json:"userID"`
}
