package hypnus

import (
	"fmt"
	logger "go-open/library/log"
	"path"
	"runtime"
)

func JoinPaths(path1, path2 string) string {
	if path1 == "" {
		return path1
	}
	path := path.Join(path1, path2)

	lastChar := func(str string) uint8 {
		if str == "" {
			panic("the str can't be empty")
		}

		return str[len(str)-1]
	}

	if lastChar(path2) == '/' && lastChar(path) != '/' {
		return path + "/"
	}
	return path
}

func Recovery() {
	if err := recover(); err != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		rs := runtime.Stack(buf, false)
		if size < rs {
			rs = size
		}
		str := fmt.Sprintf("%s", buf[:rs])
		logger.Error(str, err)
	}
}
