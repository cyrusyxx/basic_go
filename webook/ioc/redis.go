package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	// Get Redis Addr From Viper
	redisaddr := viper.GetString("redis.addr")

	// Init Redis
	redisdb := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})

	return redisdb
}
