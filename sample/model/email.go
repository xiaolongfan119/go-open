package model

import (
	"github.com/xiaolongfan119/go-open/library/database/orm"
)

type Email struct {
	orm.Model
	Email  string `json:"email"`
	UserID uint   `json:"userID"`
}
