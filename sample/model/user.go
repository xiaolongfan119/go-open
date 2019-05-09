package model

import (
	"go-open/library/database/orm"
)

type User struct {
	orm.Model
	Age   int
	Name  string
	Email string
}

// type Email struct {
// 	ID     int
// 	UserID int
// 	Email  string
// }
