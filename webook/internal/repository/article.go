package repository

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, arti domain.Article) (int64, error)
	Update(ctx context.Context, arti domain.Article) error
	Sync(ctx context.Context, arti domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDAO

	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO
}

func NewCachedArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (r *CachedArticleRepository) Create(ctx context.Context,
	arti domain.Article) (int64, error) {
	return r.dao.Insert(ctx, r.toEntity(ctx, arti))
}

func (r *CachedArticleRepository) toEntity(ctx context.Context,
	arti domain.Article) dao.Article {
	return dao.Article{
		Id:       arti.Id,
		AuthorId: arti.Author.Id,
		Title:    arti.Title,
		Content:  arti.Content,
	}
}

func (r *CachedArticleRepository) Update(ctx context.Context,
	arti domain.Article) error {
	return r.dao.UpdateById(ctx, r.toEntity(ctx, arti))
}

func (r *CachedArticleRepository) Sync(ctx context.Context,
	arti domain.Article) (int64, error) {
	return r.dao.Sync(ctx, r.toEntity(ctx, arti))
}
