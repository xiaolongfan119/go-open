package main

import (
	"github.com/xiaolongfan119/go-open/v2/sample/model"

	http "github.com/xiaolongfan119/go-open/v2/sample/http"

	"github.com/xiaolongfan119/go-open/v2/sample/conf"

	log "github.com/xiaolongfan119/go-open/v2/library/log"
)

func main() {
	conf.Init()
	log.Init(conf.Conf.LogConfig)
	model.Init(conf.Conf)
	http.Init(conf.Conf.HttpServer)
	defer func() {
		if err := recover(); err != nil {
			log.Warn(err)
		}
	}()
}
