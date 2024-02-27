package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/webook/internal/repository"
	"webook/webook/internal/service/sms"
)

var (
	ErrCodeSendTooFast   = repository.ErrCodeSendTooFast
	ErrCodeVerifyTooFast = repository.ErrCodeVerifyTooFast
)

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo *repository.CodeRepository, sms sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  sms,
	}
}

func (s *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := s.generateCode()

	err := s.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	// Send code to phone
	const tpl = "123456"
	return s.sms.Send(ctx, phone, tpl, []string{code})
}

func (s *CodeService) Verify(ctx context.Context,
	biz, phone, code string) (bool, error) {
	ok, err := s.repo.Verify(ctx, biz, phone, code)
	if err == repository.ErrCodeVerifyTooFast {
		return false, nil
	}
	return ok, err
}

func (s *CodeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
