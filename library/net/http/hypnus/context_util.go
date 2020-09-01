package hypnus

import (
	"net"
	"strconv"

	"github.com/xiaolongfan119/go-open/v2/library/ecode"
)

// 便捷方法
func (c *Context) ToInt(str string) (int, error) {

	var (
		i   int
		err error
	)

	if i, err = strconv.Atoi(str); err != nil {
		c.Abort()
		c.JSON(nil, ecode.ParamsInValid)
		return 0, err
	}
	return i, nil
}

func (c *Context) ToInts(strs ...string) ([]int, error) {

	var ints []int
	for _, str := range strs {
		if i, err := c.ToInt(str); err == nil {
			ints = append(ints, i)
		} else {
			return nil, err
		}
	}
	return ints, nil
}

func (c *Context) GetIP() string {
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	return ip
}
