package orm

import "time"

type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:createAt" json:"createAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}
