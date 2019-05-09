package log

import (
	"os"
	"sync"

	logger "github.com/sirupsen/logrus"
)

var (
	Info  = logger.Info
	Warn  = logger.Warn
	Error = logger.Error
)

var (
	packageName        string
	callerInitOnce     sync.Once
	minimumCallerDepth = 3
)

type LogConfig struct {
	Path   string
	Level  string
	OutPut string
}

func Init(conf *LogConfig) {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.InfoLevel)
	logger.SetReportCaller(true)
	logger.AddHook(newLfsHook(conf.Path))
}
