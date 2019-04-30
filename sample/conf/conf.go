package conf

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/ihornet/go-commom/library/net/http/hypnus"
)

var (
	Conf     = &Config{}
	confPath string
)

type Config struct {
	HttpServer *hypnus.ServerConf
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
