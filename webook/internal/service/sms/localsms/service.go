package localsms

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{
	}
}

func (s *Service) Send(ctx context.Context, number string,
	tplID string, args []string, numbers ...string) error {

	fmt.Println("Send SMS to", number, "with tplID", tplID, "and args", args)
	return nil
}
