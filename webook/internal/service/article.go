package service

import (
	"context"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/events/article"
	"webook/webook/internal/repository"
	"webook/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, arti domain.Article) (int64, error)
	Publish(ctx context.Context, arti domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, id int64) error
	GetByAuthor(ctx context.Context, uid int64, offset int64, limit int64) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, uid, id int64) (domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset, limit int64) ([]domain.Article, error)
}

type ImplArticleService struct {
	repo     repository.ArticleRepository
	producer article.Producer

	l logger.Logger
}

func NewImplArticleService(repo repository.ArticleRepository,
	producer article.Producer, l logger.Logger) ArticleService {
	return &ImplArticleService{
		repo:     repo,
		producer: producer,
		l:        l,
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
	uid, id int64) (domain.Article, error) {
	res, err := s.repo.GetPubById(ctx, id)

	go func() {
		if err == nil {
			er := s.producer.ProducerReadEvent(article.ReadEvent{
				Aid: id,
				Uid: uid,
			})
			if er != nil {
				s.l.Error("Failed to produce read event", logger.Error(er))
			}
		}
	}()

	return res, err
}

func (s *ImplArticleService) ListPub(ctx context.Context,
	start time.Time, offset, limit int64) ([]domain.Article, error) {
	return s.repo.ListPub(ctx, start, offset, limit)
}
