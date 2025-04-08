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
	repo     repository.CommentRepository
	userRepo repository.UserRepository
}

func NewCommentServiceImpl(repo repository.CommentRepository, userRepo repository.UserRepository) CommentService {
	return &CommentServiceImpl{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *CommentServiceImpl) Create(ctx context.Context, comment domain.Comment) (int64, error) {
	user, err := s.userRepo.FindByID(ctx, comment.User.Id)
	if err != nil {
		return 0, err
	}
	comment.User.NickName = user.NickName
	return s.repo.Create(ctx, comment)
}

func (s *CommentServiceImpl) GetByArticleId(ctx context.Context, articleId int64, offset int64, limit int64) ([]domain.Comment, error) {
	return s.repo.GetByArticleId(ctx, articleId, offset, limit)
}

func (s *CommentServiceImpl) DeleteById(ctx context.Context, id int64, userId int64) error {
	return s.repo.DeleteById(ctx, id, userId)
}
