package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Id        int64  `gorm:"primaryKey;autoIncrement"`
	Content   string `gorm:"type:text"`
	ArticleId int64  `gorm:"index"`
	UserId    int64  `gorm:"index"`
	UserName  string `gorm:"type:varchar(128)"`
	Ctime     int64
	Utime     int64
}

type CommentDAO interface {
	Insert(ctx context.Context, comment Comment) (int64, error)
	GetByArticleId(ctx context.Context, articleId int64, offset int64, limit int64) ([]Comment, error)
	DeleteById(ctx context.Context, id int64, userId int64) error
}

type GORMCommentDAO struct {
	db *gorm.DB
}

func NewGORMCommentDAO(db *gorm.DB) CommentDAO {
	return &GORMCommentDAO{
		db: db,
	}
}

func (d *GORMCommentDAO) Insert(ctx context.Context, comment Comment) (int64, error) {
	now := time.Now().UnixMilli()
	comment.Ctime = now
	comment.Utime = now
	err := d.db.WithContext(ctx).Create(&comment).Error
	return comment.Id, err
}

func (d *GORMCommentDAO) GetByArticleId(ctx context.Context, articleId int64, offset int64, limit int64) ([]Comment, error) {
	var comments []Comment
	err := d.db.WithContext(ctx).
		Where("article_id = ?", articleId).
		Offset(int(offset)).
		Limit(int(limit)).
		Order("ctime DESC").
		Find(&comments).Error
	return comments, err
}

func (d *GORMCommentDAO) DeleteById(ctx context.Context, id int64, userId int64) error {
	res := d.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userId).
		Delete(&Comment{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("评论不存在或无权删除")
	}
	return nil
}
