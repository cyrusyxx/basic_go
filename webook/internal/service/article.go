package service

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, arti domain.Article) (int64, error)
	Publish(ctx context.Context, arti domain.Article) (int64, error)
}

type ImplArticleService struct {
	repo repository.ArticleRepository
}

func NewImplArticleService(repo repository.ArticleRepository) ArticleService {
	return &ImplArticleService{
		repo: repo,
	}
}

func (s *ImplArticleService) Save(ctx context.Context, arti domain.Article) (int64, error) {
	if arti.Id > 0 {
		err := s.repo.Update(ctx, arti)
		return arti.Id, err
	} else {
		return s.repo.Create(ctx, arti)
	}
}

func (s *ImplArticleService) Publish(ctx context.Context,
	arti domain.Article) (int64, error) {
	return s.repo.Sync(ctx, arti)
}
