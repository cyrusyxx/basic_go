package startup

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: "cyrusss.top:32699",
	})
}
