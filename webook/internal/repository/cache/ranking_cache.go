package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/webook/internal/domain"
)

type RankingCache interface {
	Set(ctx context.Context, artis []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type RedisRankingCache struct {
	client     redis.Client
	key        string
	expiration time.Duration
}

func NewRedisRankingCache(client redis.Client) RankingCache {
	return &RedisRankingCache{
		client:     client,
		key:        "ranking:top_n",
		expiration: time.Minute,
	}
}

func (r *RedisRankingCache) Set(ctx context.Context, artis []domain.Article) error {
	for i := range artis {
		artis[i].Content = artis[i].Abstract()
	}
	val, err := json.Marshal(artis)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var artis []domain.Article
	err = json.Unmarshal(val, &artis)
	return artis, err
}
