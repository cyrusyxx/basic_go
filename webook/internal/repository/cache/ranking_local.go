package cache

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"time"
	"webook/webook/internal/domain"
)

type LocalRankingCache struct {
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewRankingLocalCache(expiration time.Duration) *LocalRankingCache {
	return &LocalRankingCache{}
}

func (r *LocalRankingCache) Set(ctx context.Context, artis []domain.Article) error {
	r.topN.Store(artis)
	r.ddl.Store(time.Now().Add(r.expiration))
	return nil
}

func (r *LocalRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	ddl := r.ddl.Load()
	artis := r.topN.Load()
	if len(artis) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("local cache expired")
	}
	return artis, nil
}

func (r *LocalRankingCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	artis := r.topN.Load()
	if len(artis) == 0 {
		return nil, errors.New("local cache expired")
	}
	return artis, nil
}
