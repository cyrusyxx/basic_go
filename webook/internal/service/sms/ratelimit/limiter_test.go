package ratelimit

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/webook/internal/service/sms"
	smsmocks "webook/webook/internal/service/sms/mocks"
	"webook/webook/pkg/limiter"
	limitermocks "webook/webook/pkg/limiter/mocks"
)

func TestRateLimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)

		wantErr error
	}{
		{
			name: "not in rate limit",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
			wantErr: nil,
		},
		{
			name: "be in rate limit",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return svc, l
			},
			wantErr: errLimited,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			smsService, l := tc.mock(ctrl)
			svc := NewRateLimitSMSService(smsService, l, "key")

			err := svc.Send(nil, "123",
				"123", []string{"123"}, "123")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
