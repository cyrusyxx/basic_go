package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InteractiveDAO interface {
	IncreaseViewCount(ctx context.Context, biz string, bizId int64) error
	IncreaseViewCountBatch(ctx context.Context, bizs []string, ids []int64) error
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error)
	GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, id int64) (InteractiveCount, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]InteractiveCount, error)
}

type InteractiveCount struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`

	ViewCnt    int64
	LikeCnt    int64
	CollectCnt int64

	Ctime int64
	Utime int64
}

type UserLikeBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	Uid   int64  `gorm:"uniqueIndex:biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`

	Status int64

	Ctime int64
	Utime int64
}

type UserCollectionBiz struct {
	id    int64  `gorm:"primaryKey,autoIncrement"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`
	Uid   int64  `gorm:"uniqueIndex:biz_type_id"`
	Cid   int64  `gorm:"index;uniqueIndex:biz_type_id"` // Collections ID

	Ctime int64
	Utime int64
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}

func (d *GORMInteractiveDAO) IncreaseViewCount(ctx context.Context,
	biz string, bizId int64) error {

	now := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"view_cnt": gorm.Expr("view_cnt + 1"),
			"utime":    now,
		}),
	}).Create(&InteractiveCount{
		BizId:   bizId,
		Biz:     biz,
		ViewCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

func (d *GORMInteractiveDAO) IncreaseViewCountBatch(ctx context.Context,
	bizs []string, ids []int64) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMInteractiveDAO(tx)
		for i := range bizs {
			err := txDAO.IncreaseViewCount(ctx, bizs[i], ids[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context,
	biz string, id int64, uid int64) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			BizId:  id,
			Biz:    biz,
			Status: 1,
			Utime:  now,
			Ctime:  now,
		}).Error
		if err != nil {
			return err
		}
		return d.db.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("like_cnt + 1"),
				"utime":    now,
			}),
		}).Create(&InteractiveCount{
			Biz:     biz,
			BizId:   id,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (d *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context,
	biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).
			Where("uid = ? AND biz = ? AND biz_id = ?", uid, biz, id).
			Updates(map[string]interface{}{
				"status": 0,
				"utime":  now,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&InteractiveCount{}).
			Where("biz = ? AND biz_id = ?", biz, id).
			Updates(map[string]interface{}{
				"like_cnt": gorm.Expr("like_cnt - 1"),
				"utime":    now,
			}).Error
	})
}

func (d *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context,
	biz string, id int64, cid int64, uid int64) error {
	now := time.Now().UnixMilli()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&UserCollectionBiz{
			Biz:   biz,
			BizId: id,
			Uid:   uid,
			Cid:   cid,
			Ctime: now,
			Utime: now,
		}).Error
		if err != nil {
			return err
		}
		return d.db.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("collect_cnt + 1"),
				"utime":       now,
			}),
		}).Create(&InteractiveCount{
			Biz:        biz,
			BizId:      id,
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
		}).Error
	})
}

func (d *GORMInteractiveDAO) GetLikeInfo(ctx context.Context,
	biz string, id int64, uid int64) (UserLikeBiz, error) {
	var like UserLikeBiz
	err := d.db.
		WithContext(ctx).
		Where("uid = ? AND biz = ? AND biz_id = ? AND status = ?",
			uid, biz, id, 1).
		First(&like).Error
	return like, err
}

func (d *GORMInteractiveDAO) GetCollectInfo(ctx context.Context,
	biz string, id int64, uid int64) (UserCollectionBiz, error) {
	var collect UserCollectionBiz
	err := d.db.
		WithContext(ctx).
		Where("uid = ? AND biz = ? AND biz_id = ?",
			uid, biz, id).
		First(&collect).Error
	return collect, err
}

func (d *GORMInteractiveDAO) Get(ctx context.Context, biz string, id int64) (InteractiveCount, error) {
	var count InteractiveCount
	err := d.db.
		WithContext(ctx).
		Where("biz = ? AND biz_id = ?", biz, id).
		First(&count).Error
	return count, err
}

func (d *GORMInteractiveDAO) GetByIds(ctx context.Context,
	biz string, ids []int64) ([]InteractiveCount, error) {
	var counts []InteractiveCount
	err := d.db.WithContext(ctx).
		Where("biz = ? AND biz_id IN ?", biz, ids).
		Find(&counts).Error
	return counts, err
}
