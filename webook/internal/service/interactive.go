package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

type InteractiveService interface {
	IncreaseViewCount(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.InteractiveCount, error)
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
		inter.Liked, er = s.repo.Collected(ctx, biz, id, uid)
		return er
	})

	return inter, eg.Wait()
}
