package model

import (
	"time"

	"github.com/ihornet/go-open/library/database/orm"
)

type User struct {
	orm.Model
	Age     int       `json:"age" default:"5" body:"age"`
	Name    string    `json:"name" body:"name" validate:"min=10"`
	Time    time.Time `json:"time" body:"time"`
	IsMan   bool      `json:"isMan" body:"isMan"`
	Email   Email     `json:"email"`
	Profile Profile   `json:"profile"`
}
