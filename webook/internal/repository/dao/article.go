package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Title    string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`

	Ctime int64
	Utime int64
}

type PublicArticle Article

type ArticleDAO interface {
	Insert(ctx context.Context, arti Article) (int64, error)
	UpdateById(ctx context.Context, arti Article) error
	Sync(ctx context.Context, entity Article) (int64, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (d *GORMArticleDAO) Insert(ctx context.Context, arti Article) (int64, error) {
	now := time.Now().UnixMilli()
	arti.Ctime = now
	arti.Utime = now
	err := d.db.WithContext(ctx).Create(&arti).Error
	return arti.Id, err
}

func (d *GORMArticleDAO) UpdateById(ctx context.Context, arti Article) error {
	now := time.Now().UnixMilli()
	res := d.db.WithContext(ctx).Model(&Article{}).
		Where("id =? AND author_id = ?", arti.Id, arti.AuthorId).
		Updates(map[string]any{
			"title":   arti.Title,
			"content": arti.Content,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("update article fail, " +
			"Id or Author is wrong")
	}
	return nil
}

func (d *GORMArticleDAO) Sync(ctx context.Context, arti Article) (int64, error) {
	tx := d.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()

	var (
		id  = arti.Id
		err error
	)
	dao := NewGORMArticleDAO(tx)
	if id > 0 {
		err = dao.UpdateById(ctx, arti)
	} else {
		id, err = dao.Insert(ctx, arti)
	}
	if err != nil {
		return 0, err
	}
	arti.Id = id
	pubArti := PublicArticle(arti)
	pubArti.Ctime = time.Now().UnixMilli()
	pubArti.Utime = time.Now().UnixMilli()
	err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   pubArti.Title,
			"content": pubArti.Content,
			"utime":   time.Now().UnixMilli(),
		}),
	}).Create(pubArti).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}
