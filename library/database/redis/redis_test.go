package redis

import (
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {

	config := RedisConfig{Addr: "192.168.1.152:6379"}
	client := NewRedisClient(&config)

	vvv, err := client.Get("xxx").Result()
	if err != nil {
		println("err")
		println(err)
	} else {
		println(vvv)
	}

	client.Set("platform:user:1:token", "xxxxxxxx", time.Second*2)
	client.Set("platform:user:1:token", "xxxxyyyyyxxxx", time.Second*4)

	client.SAdd("key", "stoken")

	v, e := client.SIsMember("key", "stoken").Result()
	println("==========")
	println(v, e)

	val, err := client.Get("platform:user:1:token").Result()
	println(val, err)

	client.Del("platform:user:1:token")
	val, err = client.Get("platform:user:1:token").Result()

	if err != nil {
		println(err.Error())
	}

}
