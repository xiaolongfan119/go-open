package orm

import "time"

type Model struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:createAt" json:"-"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"-"`
}
