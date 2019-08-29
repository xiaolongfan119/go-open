package orm

import (
	"fmt"
	"time"

	xtime "github.com/ihornet/go-open/library/time"

	log "github.com/ihornet/go-open/library/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

type DBConfig struct {
	Host     string
	DBName   string
	Username string
	Password string
	Port     string

	Active      int
	Idle        int
	IdleTimeout xtime.Duration
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info(v...)
}

func NewMySQL(conf *DBConfig) (db *gorm.DB) {
	var err error
	path := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName)
	DB, err = gorm.Open("mysql", path)
	if err != nil {
		log.Info(fmt.Sprintf("db connect with path: %s   err: %v \n", path, err))
		panic(err)
	}
	DB.DB().SetMaxIdleConns(conf.Idle)
	DB.DB().SetMaxOpenConns(conf.Active)
	DB.DB().SetConnMaxLifetime(time.Duration(conf.IdleTimeout) / time.Second)
	DB.SetLogger(ormLog{})
	DB.SingularTable(true)
	DB.LogMode(true)
	return DB
}
