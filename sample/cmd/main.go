package main

import (
	"github.com/xiaolongfan119/go-open/sample/model"

	http "github.com/xiaolongfan119/go-open/sample/http"

	"github.com/xiaolongfan119/go-open/sample/conf"

	log "github.com/xiaolongfan119/go-open/library/log"
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
