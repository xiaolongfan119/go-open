package model

import (
	"go-open/library/database/orm"
)

type Email struct {
	orm.Model
	Email  string `json:"email"`
	UserID uint   `json:"userID"`
}
