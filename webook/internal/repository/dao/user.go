package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Email    string
	Password string
	// Create time
	Ctime int64
	// Update time
	Utime int64
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	return dao.db.WithContext(ctx).Create(&u).Error
}
