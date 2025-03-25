package cache

import (
	"context"
	"errors"
	"time"
	"webook/webook/internal/domain"

	"github.com/ecodeclub/ekit/syncx/atomicx"
)

type LocalRankingCache struct {
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewRankingLocalCache() *LocalRankingCache {
	return &LocalRankingCache{
		topN:       atomicx.NewValue[[]domain.Article](),
		ddl:        atomicx.NewValue[time.Time](),
		expiration: time.Minute,
	}
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
