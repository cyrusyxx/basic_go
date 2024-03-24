package constants

import "time"

const (
	RateLimitRate = 100
)

var (
	RateLimitInterval    = time.Minute
	MaxCorsAge           = 12 * time.Hour
	JwtExpireTime        = 30 * time.Minute
	JwtRefreshExpireTime = 7 * 24 * time.Hour
	CheckLoginExpireTime = 30 * time.Minute
	UserCacheExpireTime  = 15 * time.Minute
)
