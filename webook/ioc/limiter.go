package ioc

import (
	"github.com/redis/go-redis/v9"
	"webook/webook/constants"
	"webook/webook/pkg/limiter"
)

func InitLimiter(rediscmd redis.Cmdable) limiter.Limiter {
	return limiter.NewRedisSlidingWindowLimiter(rediscmd,
		constants.RateLimitInterval, constants.RateLimitRate)
}
