package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/ihornet/go-open/library/log"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var RedisClient *redis.Client

func NewRedisClient(config *RedisConfig) (client *redis.Client) {

	_client := redis.NewClient(&redis.Options{Addr: config.Addr, Password: config.Password, DB: config.DB})
	if _, err := _client.Ping().Result(); err != nil {
		log.Info(fmt.Sprintf("redis connect with addr: %s   err: %v \n", config.Addr, err))
		panic(err)
	} else {
		log.Info(fmt.Sprintf("redis connect with addr: %s   success \n", config.Addr))
	}

	RedisClient = _client
	return RedisClient
}
