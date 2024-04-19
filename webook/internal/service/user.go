package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

var (
	ErrInvalidUserOrPassword = errors.New("user not found or password is wrong")
	ErrDuplicateUser         = repository.ErrDuplicateUser
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	Edit(ctx context.Context, uid int64, nickname string, birthday string, description string) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type CachedUserService struct {
	repo repository.UserRepository
}

func NewCachedUserService(repo repository.UserRepository) UserService {
	return &CachedUserService{
		repo: repo,
	}
}

func (svc *CachedUserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *CachedUserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// Find User by Email
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	// Check Password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *CachedUserService) Edit(ctx context.Context, uid int64, nickname string, birthday string, description string) error {
	return svc.repo.EditProfile(ctx, domain.User{
		Id:          uid,
		NickName:    nickname,
		Birthday:    birthday,
		Description: description,
	})
}

func (svc *CachedUserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindByID(ctx, id)
}

func (svc *CachedUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {

	// Find User by Phone
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err == repository.ErrUserNotFound {
		// Create User
		err = svc.repo.Create(ctx, domain.User{
			Phone: phone,
		})
		// If user exists, return user
		if err == repository.ErrDuplicateUser {
			return svc.repo.FindByPhone(ctx, phone)
		}
		// If error is not duplicate user, return error
		if err != nil {
			return domain.User{}, err
		}
		// If user not exists, return user
		return svc.repo.FindByPhone(ctx, phone)
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *CachedUserService) FindOrCreateByWechat(ctx context.Context,
	info domain.WechatInfo) (domain.User, error) {
	// Find User by Wechat OpenID
	u, err := svc.repo.FindByWechatOpenID(ctx, info.OpenId)
	if err == repository.ErrUserNotFound {
		// Create User
		err = svc.repo.Create(ctx, domain.User{
			WechatInfo: info,
		})

		// If user exists, return user
		if err == repository.ErrDuplicateUser {
			return svc.repo.FindByWechatOpenID(ctx, info.OpenId)
		}
		// If error is not duplicate user, return error
		if err != nil {
			return domain.User{}, err
		}
		// If user not exists, return user
		return svc.repo.FindByWechatOpenID(ctx, info.OpenId)
	}
	return u, err
}
