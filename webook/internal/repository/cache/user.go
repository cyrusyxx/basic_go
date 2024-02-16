package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/webook/constants"
	"webook/webook/internal/domain"
)

type UserCache struct {
	cmd         redis.Cmdable
	exprireTime time.Duration
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:         cmd,
		exprireTime: constants.UserCacheExpireTime,
	}
}

func (cache *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	// Get key of id
	key := cache.key(id)

	// Get data from cache
	data, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	// Unmarshal data to domain.User
	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return domain.User{}, err
	}

	// Return user
	return user, nil
}

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {
	// Marshal user to data
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Set data to cache
	key := cache.key(user.Id)
	return cache.cmd.Set(ctx, key, data, cache.exprireTime).Err()
}
