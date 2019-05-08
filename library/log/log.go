package log

import (
	"runtime"
	"strings"
	"sync"

	logger "github.com/sirupsen/logrus"
)

var (
	Info  logger.Info
	Warn  logger.Warn
	Error logger.Error
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

// func Init(conf *LogConfig) {
// 	logger.SetOutput(os.Stdout)
// 	logger.SetLevel(logger.InfoLevel)
// 	logger.SetReportCaller(true)
// 	logger.AddHook(newLfsHook(conf.Path))
// }

// func Info(args ...interface{}) {
// 	logger.Info(args...)
// }

// func Warn(args ...interface{}) {
// 	logger.Warn(args...)
// }

// func Error(args ...interface{}) {
// 	logger.Error(append(args, *getCaller()))
// }

func getCaller() *runtime.Frame {

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		packageName = getPackageName(runtime.FuncForPC(pcs[1]).Name())
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, minimumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != packageName {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
