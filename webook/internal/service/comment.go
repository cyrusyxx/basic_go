package service

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

type CommentService interface {
	Create(ctx context.Context, comment domain.Comment) (int64, error)
	GetByArticleId(ctx context.Context, articleId int64, offset int64, limit int64) ([]domain.Comment, error)
	DeleteById(ctx context.Context, id int64, userId int64) error
}

type CommentServiceImpl struct {
	repo repository.CommentRepository
}

func NewCommentServiceImpl(repo repository.CommentRepository) CommentService {
	return &CommentServiceImpl{
		repo: repo,
	}
}

func (s *CommentServiceImpl) Create(ctx context.Context, comment domain.Comment) (int64, error) {
	return s.repo.Create(ctx, comment)
}

func (s *CommentServiceImpl) GetByArticleId(ctx context.Context, articleId int64, offset int64, limit int64) ([]domain.Comment, error) {
	return s.repo.GetByArticleId(ctx, articleId, offset, limit)
}

func (s *CommentServiceImpl) DeleteById(ctx context.Context, id int64, userId int64) error {
	return s.repo.DeleteById(ctx, id, userId)
}
