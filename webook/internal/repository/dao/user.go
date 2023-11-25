package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	Id       int64  `gorm:"primaryKey, autoncrement"` // There is a Bug!!!
	Email    string `gorm:"unique"`
	Password string
	// Create time
	Ctime int64
	// Update time
	Utime int64
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}

}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}
