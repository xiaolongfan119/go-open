package main

import (
	log "go-open/library/log"
	"go-open/sample/conf"
	http "go-open/sample/http"
	"go-open/sample/model"
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
