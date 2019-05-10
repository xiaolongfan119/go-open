package log

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

var (
	Info = logger.Info
	Warn = logger.Warn
)

func Error(stack interface{}, args interface{}) {
	if stack == nil {
		logger.Error(args)
	} else {
		logger.WithFields(logger.Fields{"msg": args}).Error(stack)
	}
}

type LogConfig struct {
	Path   string
	Level  string
	OutPut string
}

func Init(conf *LogConfig) {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.InfoLevel)
	// logger.SetReportCaller(true)
	logger.AddHook(newLfsHook(conf.Path))
}
