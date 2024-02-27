package repository

import (
	"context"
	"webook/webook/internal/repository/cache"
)

var (
	ErrCodeVerifyTooFast = cache.ErrCodeVerifyTooFast
)

type CodeRepository struct {
	cache cache.CodeCache
}

func (c *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
