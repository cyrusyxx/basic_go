package constants

import "time"

const (
	RateLimitRate = 100
)

var (
	RateLimitInterval = time.Minute
	MaxCorsAge        = 12 * time.Hour
	//JwtExpireTime          = 30 * time.Minute
	JwtExpireTime          = 100 * time.Hour
	JwtRefreshExpireTime   = 7 * 24 * time.Hour
	CheckLoginExpireTime   = 30 * time.Minute
	UserCacheExpireTime    = 1 * time.Minute
	InteractiveCacheExpire = 1 * time.Minute
)
