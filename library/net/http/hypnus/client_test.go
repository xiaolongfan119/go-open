package hypnus

import (
	"context"
	"fmt"
	breaker "go-open/library/net/netutil/breaker"
	xtime "go-open/library/time"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {

	c := &ClientConfig{
		TimeOut:   xtime.Duration(time.Duration(10) * time.Second),
		KeepAlive: xtime.Duration(time.Duration(10) * time.Second),
		bkConfig:  &breaker.BreakerConfig{Timeout: xtime.Duration(time.Duration(10) * time.Second)},
	}

	client := NewClient(c)

	var isFail = false
	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		//	go func() {
		var resp Resp
		if err := client.Get(context.TODO(), "http://localhost:8081/", &resp); err != nil {
			isFail = true
			fmt.Println("############ isfail ")
		} else {
			fmt.Println("------------ success ")
		}
		wg.Done()

		if i == 57 {
			time.Sleep(time.Duration(8) * time.Second)
			fmt.Println("-+++++++++++++ ")
		}

		if i == 60 {
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Println("-+++++++++++++ ")
		}

		//}()
	}

	wg.Wait()
	assert.Equal(t, isFail, false)
}

type Resp struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		Id         int    `json:"id"`
		Name       string `json:"name"`
		Parent     string `json:"parent"`
		Code       string `json:"code"`
		Level      string `json:"level"`
		Mark       string `json:"mark"`
		ParentCode string `json:"parentCode"`
	} `json:"data"`
}
