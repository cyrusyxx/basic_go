package failover

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/webook/internal/service/sms"
	smsmocks "webook/webook/internal/service/sms/mocks"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name  string
		mocks func(ctrl *gomock.Controller) []sms.Service

		wantErr error
	}{
		{
			name: "success at the first time",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			wantErr: nil,
		},
		{
			name: "success at the second time",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("send fail"))
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			wantErr: nil,
		},
		{
			name: "all fail",
			mocks: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("send fail"))
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("send fail"))
				return []sms.Service{svc0, svc1}
			},
			wantErr: errors.New("all services failed to send sms"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewFailOverSMSService(tc.mocks(ctrl)...)
			err := svc.Send(nil, "", "", nil)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
