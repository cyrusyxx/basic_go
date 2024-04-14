package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/webook/internal/domain"
	"webook/webook/internal/service"
	svcmocks "webook/webook/internal/service/mocks"
	"webook/webook/internal/web/jwt"
	jwtmocks "webook/webook/internal/web/jwt/mocks"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService,
			service.CodeService, jwt.Handler)

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody string
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, jwt.Handler) {
				userService := svcmocks.NewMockUserService(ctrl)
				userService.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "12345@qq.com",
					Password: "hello#world123",
				}).Return(nil)

				codeService := svcmocks.NewMockCodeService(ctrl)
				jwthandler := jwtmocks.NewMockHandler(ctrl)

				return userService, codeService, jwthandler
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email":"12345@qq.com",
"password":"hello#world123",
"confirmPassword":"hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantBody: "HELLO ITS IN SIGNUP",
		},

		{
			name: "Email Pattern is wrong",
			mock: func(ctrl *gomock.Controller) (service.UserService,
				service.CodeService, jwt.Handler) {
				userService := svcmocks.NewMockUserService(ctrl)
				codeService := svcmocks.NewMockCodeService(ctrl)
				jwthandler := jwtmocks.NewMockHandler(ctrl)

				return userService, codeService, jwthandler
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email":"12345",
"password":"hello#world123",
"confirmPassword":"hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantBody: "Email pattern is wrong\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService, codeService, jwthandler := tc.mock(ctrl)
			hdl := NewUserHandler(userService, codeService, jwthandler)

			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}
