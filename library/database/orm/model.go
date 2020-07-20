package orm

import "time"

type Model struct {
	ID             int       `gorm:"primary_key" json:"id"`
	CreateTime     time.Time `gorm:"column:createTime" json:"-"`
	LastUpdateTime time.Time `gorm:"column:lastUpdateTime" json:"-"`
}
