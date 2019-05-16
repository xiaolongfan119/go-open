package model

import (
	"go-open/library/database/orm"
)

type User struct {
	orm.Model
	Age     int     `json:"age" default:"5"`
	Name    string  `json:"name"`
	Email   Email   `json:"email"`
	Profile Profile `json:"profile"`
}
