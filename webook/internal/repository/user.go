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

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	EditProfile(ctx context.Context, u domain.User) error
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
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

func (repo *CachedUserRepository) EditProfile(ctx context.Context, u domain.User) error {
	return repo.dao.Edit(ctx, dao.User{
		Id:          u.Id,
		NickName:    u.NickName,
		Birthday:    u.Birthday,
		Description: u.Description,
	})
}

func (repo *CachedUserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
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

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
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

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
