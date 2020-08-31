package conf

import (
	"flag"

	"github.com/xiaolongfan119/go-open/library/database/orm"

	"github.com/xiaolongfan119/go-open/library/net/http/hypnus"

	jwt "github.com/xiaolongfan119/go-open/library/net/http/hypnus/middleware/token"

	log "github.com/xiaolongfan119/go-open/library/log"

	"github.com/BurntSushi/toml"
)

var (
	Conf     = &Config{}
	confPath string
)

type Config struct {
	HttpServer  *hypnus.ServerConf
	TokenConfig *jwt.TokenConfig
	LogConfig   *log.LogConfig
	DBConfig    *orm.DBConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

func Init() (err error) {
	confPath = "./../cmd/test-conf.toml"
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func remote() (err error) {
	// TODO
	return
}
