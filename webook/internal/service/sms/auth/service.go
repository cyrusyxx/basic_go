package auth

import (
	"context"
	"webook/webook/internal/service/sms"
)

type AuthSMSService struct {
	svc sms.Service
}

func (s *AuthSMSService) Send(ctx context.Context, number string, tplToken string, args []string, numbers ...string) error {
	//TODO implement me
	//TODO in 五.12.
	return nil
}
