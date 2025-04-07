package repository

import (
	"context"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/dao"
)

type CommentRepository interface {
	Create(ctx context.Context, comment domain.Comment) (int64, error)
	GetByArticleId(ctx context.Context, articleId int64, offset int64, limit int64) ([]domain.Comment, error)
	DeleteById(ctx context.Context, id int64, userId int64) error
}

type CommentRepo struct {
	dao dao.CommentDAO
}

func NewCommentRepo(dao dao.CommentDAO) CommentRepository {
	return &CommentRepo{
		dao: dao,
	}
}

func (r *CommentRepo) Create(ctx context.Context, comment domain.Comment) (int64, error) {
	return r.dao.Insert(ctx, dao.Comment{
		Content:   comment.Content,
		ArticleId: comment.ArticleId,
		UserId:    comment.User.Id,
		UserName:  comment.User.NickName,
	})
}

func (r *CommentRepo) GetByArticleId(ctx context.Context,
	articleId int64, offset int64, limit int64) ([]domain.Comment, error) {

	comments, err := r.dao.GetByArticleId(ctx, articleId, offset, limit)
	if err != nil {
		return nil, err
	}

	return r.toDomain(comments), nil
}

func (r *CommentRepo) DeleteById(ctx context.Context, id int64, userId int64) error {
	return r.dao.DeleteById(ctx, id, userId)
}

func (r *CommentRepo) toDomain(comments []dao.Comment) []domain.Comment {
	res := make([]domain.Comment, 0, len(comments))
	for _, comment := range comments {
		res = append(res, domain.Comment{
			Id:        comment.Id,
			Content:   comment.Content,
			ArticleId: comment.ArticleId,
			User: domain.User{
				Id:       comment.UserId,
				NickName: comment.UserName,
			},
			Ctime: time.UnixMilli(comment.Ctime),
			Utime: time.UnixMilli(comment.Utime),
		})
	}
	return res
}
