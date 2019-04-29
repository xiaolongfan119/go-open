package hypnus

import "path"

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
