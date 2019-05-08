package log

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func newLfsHook(path string) log.Hook {

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.InfoLevel:  getWriter(log.InfoLevel, path),
		log.WarnLevel:  getWriter(log.WarnLevel, path),
		log.ErrorLevel: getWriter(log.ErrorLevel, path),
	}, &log.TextFormatter{DisableColors: false})

	return lfsHook
}

func getWriter(level log.Level, path string) (w *rotatelogs.RotateLogs) {
	var (
		fileName      string
		rotationTime  int
		rotationCount uint
	)

	switch {
	case level == log.InfoLevel:
		fileName = "log"
		rotationTime = 24
		rotationCount = 30
	case level == log.WarnLevel:
		fileName = "warn"
		rotationTime = 7
		rotationCount = 5
	case level == log.ErrorLevel:
		fileName = "error"
		rotationTime = 15
		rotationCount = 5
	}

	w, err := rotatelogs.New(
		path+fileName+".%Y%m%d",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(fileName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(rotationCount),
	)

	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}
	return
}
