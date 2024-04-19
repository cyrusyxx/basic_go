package service

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, arti domain.Article) (int64, error)
	Publish(ctx context.Context, arti domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, id int64) error
	GetByAuthor(ctx context.Context, uid int64, offset int64, limit int64) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
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
	arti.Status = domain.ArticleStatusUnpublished
	if arti.Id > 0 {
		err := s.repo.Update(ctx, arti)
		return arti.Id, err
	} else {
		return s.repo.Create(ctx, arti)
	}
}

func (s *ImplArticleService) Publish(ctx context.Context,
	arti domain.Article) (int64, error) {
	arti.Status = domain.ArticleStatusPublished
	return s.repo.Sync(ctx, arti)
}

func (s *ImplArticleService) Withdraw(ctx context.Context,
	uid int64, id int64) error {
	return s.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)
}

func (s *ImplArticleService) GetByAuthor(ctx context.Context,
	uid int64, offset int64, limit int64) ([]domain.Article, error) {
	return s.repo.GetByAuthor(ctx, uid, offset, limit)
}

func (s *ImplArticleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return s.repo.GetById(ctx, id)
}

func (s *ImplArticleService) GetPubById(ctx context.Context,
	id int64) (domain.Article, error) {
	return s.repo.GetPubById(ctx, id)
}
