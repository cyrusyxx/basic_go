package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

var (
	ErrInvalidUserOrPassword = errors.New("user not found or password is wrong")
	ErrDuplicateUser         = repository.ErrDuplicateUser
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
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

func (svc *UserService) Edit(ctx context.Context, uid int64, nickname string, birthday string, description string) error {
	return svc.repo.EditProfile(ctx, domain.User{
		Id:          uid,
		NickName:    nickname,
		Birthday:    birthday,
		Description: description,
	})
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindByID(ctx, id)
}

func (s *UserService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	// Find User by Phone
	u, err := s.repo.FindByPhone(ctx, phone)
	if err == repository.ErrUserNotFound {
		// TODO: Finish repo.Create
		// Create User
		err = s.repo.Create(ctx, domain.User{
			Phone: phone,
		})
		// If user exists, return user
		if err == repository.ErrDuplicateUser {
			return s.repo.FindByPhone(ctx, phone)
		}
		// If error is not duplicate user, return error
		if err != nil {
			return domain.User{}, err
		}
		// If user not exists, return user
		return s.repo.FindByPhone(ctx, phone)
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
