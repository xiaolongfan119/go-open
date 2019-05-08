package orm

import (
	"fmt"
	xtime "go-open/library/time"
	"time"

	log "go-open/library/log"

	"github.com/jinzhu/gorm"
)

type Config struct {
	Path        string
	Active      int
	Idle        int
	IdleTimeout xtime.Duration
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(v...)
}

func NewMySQL(conf *Config) (db *gorm.DB) {
	db, err := gorm.Open("mysql", conf.Path)
	if err != nil {
		log.Info(fmt.Sprintf("db connect with path: %s   err: %v \n", conf.Path, err))
		panic(err)
	}
	db.DB().SetMaxIdleConns(conf.Idle)
	db.DB().SetMaxOpenConns(conf.Active)
	db.DB().SetConnMaxLifetime(time.Duration(conf.IdleTimeout) / time.Second)
	db.SetLogger(ormLog{})
	return
}
