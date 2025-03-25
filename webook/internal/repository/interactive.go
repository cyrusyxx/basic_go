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
	IncreaseViewCountBatch(ctx context.Context, bizs []string, bizIds []int64) error
	IncreaseLike(ctx context.Context, biz string, id int64, uid int64) error
	DecreaseLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.InteractiveCount, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.InteractiveCount, error)
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

func (r *CachedInteractiveRepository) IncreaseViewCountBatch(ctx context.Context,
	bizs []string, bizIds []int64) error {
	err := r.dao.IncreaseViewCountBatch(ctx, bizs, bizIds)
	if err != nil {
		return err
	}
	go func() {
		for i := range bizs {
			er := r.cache.IncreaseViewCountIfPresent(ctx, bizs[i], bizIds[i])
			if er != nil {
				r.l.Error("Failed to increase view count Cache",
					logger.Error(er),
					logger.String("biz", bizs[i]),
					logger.Int64("bizId", bizIds[i]),
				)
			}
		}
	}()
	return nil
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

	c, err := r.cache.Get(ctx, biz, id)
	if err == nil {
		return c, nil
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

func (r *CachedInteractiveRepository) GetByIds(ctx context.Context,
	biz string, ids []int64) ([]domain.InteractiveCount, error) {
	inters, err := r.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	return r.toDomains(inters), nil
}

func (r *CachedInteractiveRepository) toDomain(inter dao.InteractiveCount) domain.InteractiveCount {
	return domain.InteractiveCount{
		BizId:      inter.BizId,
		ViewCnt:    inter.ViewCnt,
		LikeCnt:    inter.LikeCnt,
		CollectCnt: inter.CollectCnt,
	}
}

// toDomains
func (r *CachedInteractiveRepository) toDomains(inters []dao.InteractiveCount) []domain.InteractiveCount {
	res := make([]domain.InteractiveCount, len(inters))
	for i := range inters {
		res[i] = r.toDomain(inters[i])
	}
	return res
}
