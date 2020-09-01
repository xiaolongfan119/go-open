package hypnus

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strconv"

	logger "github.com/xiaolongfan119/go-open/v2/library/log"
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

func convertMap2StrMap(data map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k := range data {
		switch data[k].(type) {
		case string:
			m[k] = data[k].(string)
		case bool:
			m[k] = strconv.FormatBool(data[k].(bool))
		case float64:
			m[k] = strconv.FormatFloat(data[k].(float64), 'f', -1, 64)
		default:
			fmt.Printf("%v", reflect.TypeOf(data[k]))
		}
	}
	return m
}
