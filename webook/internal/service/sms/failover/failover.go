package failover

import (
	"context"
	"errors"
	"log"
	"webook/webook/internal/service/sms"
)

type FailOverSMSService struct {
	svcs []sms.Service
}

func NewFailOverSMSService(svcs ...sms.Service) *FailOverSMSService {
	return &FailOverSMSService{svcs: svcs}
}

func (f *FailOverSMSService) Send(ctx context.Context, number string,
	tplID string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, number, tplID, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("all services failed to send sms")
}

func (f *FailOverSMSService) SendWithIdx(ctx context.Context, idx int, number string,
	tplID string, args []string, numbers ...string) error {
	// TODO: 负载均衡式自动切换服务商
	// TODO: in 五.11.
	return errors.New("none")
}
