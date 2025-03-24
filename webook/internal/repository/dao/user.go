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

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	Edit(ctx context.Context, u User) error
	FindByWechatOpenID(ctx context.Context, openId string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

type User struct {
	Id       int64          `gorm:"primaryKey;autoIncrement"`
	Email    sql.NullString `gorm:"type:varchar(100);unique;comment:邮箱"`
	Password string         `gorm:"type:varchar(255);not null;comment:密码"`
	Phone    sql.NullString `gorm:"type:varchar(100);unique;comment:手机号"`

	NickName    string `gorm:"type:varchar(100);comment:昵称"`
	Birthday    string `gorm:"type:varchar(50);comment:生日"`
	Description string `gorm:"type:varchar(1000);comment:个人简介"`

	WechatOpenId  sql.NullString `gorm:"type:varchar(100);unique;comment:微信开放ID"`
	WechatUnionId sql.NullString `gorm:"type:varchar(100);comment:微信联合ID"`

	Ctime int64 `gorm:"type:bigint;comment:创建时间"`
	Utime int64 `gorm:"type:bigint;comment:更新时间"`
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Edit(ctx context.Context, u User) error {
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

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
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

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByWechatOpenID(ctx context.Context,
	openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).
		Where("wechat_open_id=?", openId).First(&u).Error
	return u, err
}
