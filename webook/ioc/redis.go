package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
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

func InitRlockClient(client redis.Cmdable) *rlock.Client {
	return rlock.NewClient(client)
}
