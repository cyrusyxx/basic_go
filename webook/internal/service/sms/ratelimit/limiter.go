package ratelimit

import (
	"context"
	"errors"
	"webook/webook/internal/service/sms"
	"webook/webook/pkg/limiter"
)

var (
	errLimited = errors.New("rate limit exceeded")
)

type RateLimitSMSService struct {
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

func NewRateLimitSMSService(
	svc sms.Service, limiter limiter.Limiter, key string,
) *RateLimitSMSService {

	return &RateLimitSMSService{
		svc:     svc,
		limiter: limiter,
		key:     key,
	}
}

func (r *RateLimitSMSService) Send(
	ctx context.Context, number string,
	tplID string, args []string, numbers ...string,
) error {

	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, number, tplID, args, numbers...)
}
