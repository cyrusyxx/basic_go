package cache

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
	"webook/webook/constants"
	"webook/webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

type InteractiveCache interface {
	IncreaseViewCountIfPresent(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeIfPresent(ctx context.Context, biz string, id int64) error
	DecreaseLikeIfPresent(ctx context.Context, biz string, id int64) error
	IncreaseCollectCntIfPresent(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, id int64) (domain.InteractiveCount, error)
	Set(ctx context.Context, biz string, id int64, res domain.InteractiveCount) error
}

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldViewCnt    = "view_cnt"
	fieldLikeCnt    = "like_cnt"
	fieldCollectCnt = "collect_cnt"
)

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func NewRedisInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		client: client,
	}
}

func (c *RedisInteractiveCache) IncreaseViewCountIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := key(biz, bizId)
	return c.client.Eval(ctx, luaIncrCnt,
		[]string{key}, fieldViewCnt, +1).Err()
}

func (c *RedisInteractiveCache) IncreaseLikeIfPresent(ctx context.Context,
	biz string, id int64) error {
	key := key(biz, id)
	return c.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, +1).Err()
}

func (c *RedisInteractiveCache) DecreaseLikeIfPresent(ctx context.Context,
	biz string, id int64) error {
	key := key(biz, id)
	return c.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (c *RedisInteractiveCache) IncreaseCollectCntIfPresent(ctx context.Context,
	biz string, id int64) error {
	key := key(biz, id)
	return c.client.Eval(ctx, luaIncrCnt, []string{key}, fieldCollectCnt, +1).Err()
}

func (c *RedisInteractiveCache) Get(ctx context.Context,
	biz string, id int64) (domain.InteractiveCount, error) {
	res, err := c.client.HGetAll(ctx, key(biz, id)).Result()
	if err != nil {
		return domain.InteractiveCount{}, err
	}

	var inter domain.InteractiveCount
	inter.ViewCnt, _ = strconv.ParseInt(res[fieldViewCnt], 10, 64)
	inter.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	inter.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)

	return inter, nil
}

func (c *RedisInteractiveCache) Set(ctx context.Context,
	biz string, id int64, res domain.InteractiveCount) error {
	err := c.client.HSet(ctx, key(biz, id),
		map[string]interface{}{
			fieldViewCnt:    res.ViewCnt,
			fieldLikeCnt:    res.LikeCnt,
			fieldCollectCnt: res.CollectCnt,
		}).Err()
	if err != nil {
		return err
	}
	return c.client.Expire(ctx, key(biz, id), constants.InteractiveCacheExpire).
		Err()
}

// key func
func key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
