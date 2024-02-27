package repository

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"log"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFound  = dao.ErrRecordNotFound
	ErrKeyNotExist   = redis.Nil
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
	})
}

func (repo *UserRepository) EditProfile(ctx context.Context, u domain.User) error {
	return repo.dao.Edit(ctx, dao.User{
		Id:          u.Id,
		NickName:    u.NickName,
		Birthday:    u.Birthday,
		Description: u.Description,
	})
}

func (repo *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	// Find id from cache
	cu, err := repo.cache.Get(ctx, id)
	switch err {
	case nil:
		return cu, nil
	case ErrKeyNotExist:
		// Find id from dao
		u, err := repo.dao.FindByID(ctx, id)
		if err != nil {
			return domain.User{}, err
		}
		du := repo.toDomain(u)

		// Set cache
		err = repo.cache.Set(ctx, du)
		if err != nil {
			log.Println(err)
		}

		return du, nil
	default:
		return domain.User{}, err
	}
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:          u.Id,
		Email:       u.Email.String,
		Password:    u.Password,
		Phone:       u.Phone.String,
		NickName:    u.NickName,
		Birthday:    u.Birthday,
		Description: u.Description,
	}
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
