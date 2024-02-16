package constants

import "time"

const (
	RateLimitRate = 100
)

var (
	RateLimitInterval    = time.Minute
	MaxCorsAge           = 12 * time.Hour
	JwtExpireTime        = 30 * time.Minute
	CheckLoginExpireTime = 30 * time.Minute
)
