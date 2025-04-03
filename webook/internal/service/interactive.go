package service

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"

	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncreaseViewCount(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	CancelCollect(ctx context.Context, biz string, id int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.InteractiveCount, error)
	GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.InteractiveCount, error)
}

type ImplInteractiveService struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &ImplInteractiveService{
		repo: repo,
	}
}

func (s *ImplInteractiveService) IncreaseViewCount(ctx context.Context, biz string, bizId int64) error {
	return s.repo.IncreaseViewCount(ctx, biz, bizId)
}

func (s *ImplInteractiveService) Like(ctx context.Context,
	biz string, id int64, uid int64) error {
	return s.repo.IncreaseLike(ctx, biz, id, uid)
}

func (s *ImplInteractiveService) CancelLike(ctx context.Context,
	biz string, id int64, uid int64) error {
	return s.repo.DecreaseLike(ctx, biz, id, uid)
}

func (s *ImplInteractiveService) Collect(ctx context.Context,
	biz string, id int64, cid int64, uid int64) error {
	return s.repo.AddCollectionItem(ctx, biz, id, cid, uid)
}

func (s *ImplInteractiveService) CancelCollect(ctx context.Context,
	biz string, id int64, uid int64) error {
	return s.repo.DeleteCollectionItem(ctx, biz, id, uid)
}

func (s *ImplInteractiveService) Get(ctx context.Context,
	biz string, id int64, uid int64) (domain.InteractiveCount, error) {

	inter, err := s.repo.Get(ctx, biz, id)
	if err != nil {
		return domain.InteractiveCount{}, err
	}

	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		inter.Liked, er = s.repo.Liked(ctx, biz, id, uid)
		return er
	})

	eg.Go(func() error {
		var er error
		inter.Collected, er = s.repo.Collected(ctx, biz, id, uid)
		return er
	})

	return inter, eg.Wait()
}

func (s *ImplInteractiveService) GetByIds(ctx context.Context,
	biz string, ids []int64) (map[int64]domain.InteractiveCount, error) {

	inters, err := s.repo.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}

	// map[articleId]InteractiveCount
	res := make(map[int64]domain.InteractiveCount, len(inters))
	for _, inter := range inters {
		res[inter.BizId] = inter
	}

	return res, nil
}
