package hypnus

import (
	"context"
	"fmt"
	xtime "go-open/library/time"
	"testing"
	"time"
)

func TestGet(t *testing.T) {

	config := ClientConfig{TimeOut: xtime.Duration(time.Duration(20) * time.Second), KeepAlive: xtime.Duration(time.Duration(1) * time.Hour)}

	client := NewClient(&config)

	client.Get(context.TODO(), "http://127.0.0.1:9005/catalog2s/")

	data, err := client.Post(context.TODO(), "http://106.12.78.150:9002/users/login", map[string]interface{}{"userName": "1234", "password": "1234"})
	fmt.Println(err)
	fmt.Println(data)

}
