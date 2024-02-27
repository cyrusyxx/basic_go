package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateUser  = errors.New("user is duplicate")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	Id       int64          `gorm:"primaryKey, autoincrement"` // There is a Bug!!!
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`

	NickName    string
	Birthday    string
	Description string

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

func (dao UserDAO) Edit(ctx context.Context, u User) error {
	var fu User
	err := dao.db.WithContext(ctx).Where("id=?", u.Id).First(&fu).Error
	if err != nil {
		return err
	}

	fu.NickName = u.NickName
	fu.Birthday = u.Birthday
	fu.Description = u.Description

	err = dao.db.Save(fu).Error
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrDuplicateUser
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}

func (d *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	return u, err
}
