package repository

import (
	"context"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
