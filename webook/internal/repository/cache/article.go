package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/webook/internal/domain"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, artis []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	GetById(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, arti domain.Article) error
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, arti domain.Article) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
	}
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	key := r.firstKey(uid)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context,
	uid int64, artis []domain.Article) error {
	// Change content to abstract
	for i := 0; i < len(artis); i++ {
		artis[i].Content = artis[i].Abstract()
	}

	// Save in cache
	key := r.firstKey(uid)
	val, err := json.Marshal(artis)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, val, 10*time.Minute).Err()
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, uid int64) error {
	key := r.firstKey(uid)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisArticleCache) GetById(ctx context.Context,
	id int64) (domain.Article, error) {
	val, err := r.client.Get(ctx, r.detailKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (r *RedisArticleCache) Set(ctx context.Context,
	arti domain.Article) error {
	val, err := json.Marshal(arti)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.detailKey(arti.Id), val, 10*time.Minute).Err()
}

func (r *RedisArticleCache) GetPubById(ctx context.Context,
	id int64) (domain.Article, error) {
	val, err := r.client.Get(ctx, r.pubKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (r *RedisArticleCache) SetPub(ctx context.Context, arti domain.Article) error {
	val, err := json.Marshal(arti)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.pubKey(arti.Id), val, 10*time.Minute).Err()
}

func (r *RedisArticleCache) pubKey(id int64) string {
	return fmt.Sprintf("article:pub:detail:%d", id)
}

func (r *RedisArticleCache) detailKey(id int64) string {
	return fmt.Sprintf("article:detail:%d", id)
}

// firstKey func
func (r *RedisArticleCache) firstKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
