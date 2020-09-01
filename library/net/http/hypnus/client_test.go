package hypnus

import (
	"context"
	"fmt"
	"testing"
	"time"

	xtime "github.com/xiaolongfan119/go-open/v2/library/time"
)

func TestGet(t *testing.T) {

	config := ClientConfig{TimeOut: xtime.Duration(time.Duration(20) * time.Second), KeepAlive: xtime.Duration(time.Duration(1) * time.Hour)}

	client := NewClient(&config)

	client.Get(context.TODO(), "http://127.0.0.1:9002/experts/")

	for i := 0; i < 10001; i++ {
		data, err := client.Post(context.TODO(), "http://127.0.0.1:9002/users/login", map[string]interface{}{"userName": "admin", "password": "hys_1234"})
		fmt.Println(err)
		fmt.Println(data)
		if i == 10000 {
			println("===end")
		}
	}
	println("===end")
}
