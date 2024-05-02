package repository

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
	"webook/webook/pkg/logger"
)

type InteractiveRepository interface {
	IncreaseViewCount(ctx context.Context, biz string, bizId int64) error
	IncreaseLike(ctx context.Context, biz string, id int64, uid int64) error
	DecreaseLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.InteractiveCount, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.Logger
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO,
	cache cache.InteractiveCache, l logger.Logger) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (r *CachedInteractiveRepository) IncreaseViewCount(ctx context.Context,
	biz string, bizId int64) error {
	err := r.dao.IncreaseViewCount(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return r.cache.IncreaseViewCountIfPresent(ctx, biz, bizId)
}

func (r *CachedInteractiveRepository) IncreaseLike(ctx context.Context,
	biz string, id int64, uid int64) error {
	err := r.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return r.cache.IncreaseLikeIfPresent(ctx, biz, id)
}

func (r *CachedInteractiveRepository) DecreaseLike(ctx context.Context,
	biz string, id int64, uid int64) error {
	err := r.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return r.cache.DecreaseLikeIfPresent(ctx, biz, id)
}

func (r *CachedInteractiveRepository) AddCollectionItem(ctx context.Context,
	biz string, id int64, cid int64, uid int64) error {
	err := r.dao.InsertCollectionBiz(ctx, biz, id, cid, uid)
	if err != nil {
		return err
	}
	return r.cache.IncreaseCollectCntIfPresent(ctx, biz, id)

}

func (r *CachedInteractiveRepository) Get(ctx context.Context,
	biz string, id int64) (domain.InteractiveCount, error) {
	cache, err := r.cache.Get(ctx, biz, id)
	if err == nil {
		return cache, nil
	}

	inter, err := r.dao.Get(ctx, biz, id)
	if err != nil {
		return domain.InteractiveCount{}, err
	}

	// if no err
	res := r.toDomain(inter)
	err = r.cache.Set(ctx, biz, id, res)
	if err != nil {
		r.l.Error("failed to set cache", logger.Error(err))
	}
	return res, nil

}

func (r *CachedInteractiveRepository) toDomain(inter dao.InteractiveCount) domain.InteractiveCount {
	return domain.InteractiveCount{
		ViewCnt:    inter.ViewCnt,
		LikeCnt:    inter.LikeCnt,
		CollectCnt: inter.CollectCnt,
	}
}

func (r *CachedInteractiveRepository) Liked(ctx context.Context,
	biz string, id int64, uid int64) (bool, error) {
	_, err := r.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (r *CachedInteractiveRepository) Collected(ctx context.Context,
	biz string, id int64, uid int64) (bool, error) {
	_, err := r.dao.GetCollectInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}
