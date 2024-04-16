package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	States   uint8  `bson:"states,omitempty"`

	Ctime int64 `bson:"ctime,omitempty"`
	Utime int64 `bson:"utime,omitempty"`
}

type PublicArticle Article

type ArticleDAO interface {
	Insert(ctx context.Context, arti Article) (int64, error)
	UpdateById(ctx context.Context, arti Article) error
	Sync(ctx context.Context, entity Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error
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
			"status":  arti.States})
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
			"status":  pubArti.States,
		}),
	}).Create(pubArti).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (d *GORMArticleDAO) SyncStatus(ctx context.Context,
	uid int64, id int64, status uint8) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? AND author_id = ?", id, uid).
			Updates(map[string]any{
				"utime":  time.Now().UnixMilli(),
				"status": status,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("id or author is wrong")
		}
		return tx.Model(&PublicArticle{}).
			Where("id = ?", uid).
			Updates(map[string]any{
				"utime":  time.Now().UnixMilli(),
				"status": status,
			}).Error
	})
}
