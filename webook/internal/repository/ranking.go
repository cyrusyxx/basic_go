package repository

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, artis []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	cache cache.RankingCache

	localCache *cache.LocalRankingCache
	redisCache *cache.RedisRankingCache
}

func NewCachedRankingRepository(cache cache.RankingCache,
	localCache *cache.LocalRankingCache, redisCache *cache.RedisRankingCache) RankingRepository {
	return &CachedRankingRepository{
		cache:      cache,
		localCache: localCache,
		redisCache: redisCache,
	}
}

func (r *CachedRankingRepository) ReplaceTopN(ctx context.Context, artis []domain.Article) error {
	_ = r.localCache.Set(ctx, artis)
	return r.redisCache.Set(ctx, artis)
}

func (r *CachedRankingRepository) _unuse_GetTopN(ctx context.Context) ([]domain.Article, error) {
	return r.cache.Get(ctx)
}

func (r *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	res, err := r.localCache.Get(ctx)
	if err == nil {
		return res, nil
	}
	res, err = r.redisCache.Get(ctx)
	if err != nil {
		return r.localCache.ForceGet(ctx)
	} else {
		r.localCache.Set(ctx, res)
		return res, nil
	}
}
