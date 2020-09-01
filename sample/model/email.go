package model

import (
	"github.com/xiaolongfan119/go-open/v2/library/database/orm"
)

type Email struct {
	orm.Model
	Email  string `json:"email"`
	UserID uint   `json:"userID"`
}
