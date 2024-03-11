package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
	repomocks "webook/webook/internal/repository/mocks"
)

func TestCachedUserService_Login(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(),
					"9347553@qq.com",
				).Return(domain.User{
					Id:       1,
					Email:    "9347553@qq.com",
					Password: "$2a$10$Tc/LJO8L9ZMgvqKAfHlgOeojI1PQa.0mD3A9AU4/rVxKvgk/RevSG",
				}, nil)
				return repo
			},
			email:    "9347553@qq.com",
			password: "Cc@002300",
			wantUser: domain.User{
				Id:       1,
				Email:    "9347553@qq.com",
				Password: "$2a$10$Tc/LJO8L9ZMgvqKAfHlgOeojI1PQa.0mD3A9AU4/rVxKvgk/RevSG",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewCachedUserService(repo)

			gotUser, gotErr := svc.Login(tc.ctx, tc.email, tc.password)

			assert.Equal(t, tc.wantUser, gotUser)
			assert.Equal(t, tc.wantErr, gotErr)
		})
	}
}
