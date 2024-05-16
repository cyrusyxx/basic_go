package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webook/webook/internal/integration/startup"
	"webook/webook/pkg/ginx"
)

func TestUserHandler_SendSMSCode(t *testing.T) {
	redisdb := startup.InitRedis()
	server := startup.InitWebServer()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		wantCode int
		wnatBody ginx.Result
	}{
		{
			name: "send sms code success",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				key := "phone_code:login:18512345678"

				code, err := redisdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				duration, err := redisdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, duration > 9*time.Minute+40*time.Second)
				err = redisdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "18512345678",
			wantCode: http.StatusOK,
			wnatBody: ginx.Result{
				Msg: "Send SMS code success",
			},
		},
		{
			name: "Code Send Too Fast",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				key := "phone_code:login:18512345678"

				err := redisdb.Set(ctx, key, "123456", 10*time.Minute).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				key := "phone_code:login:18512345678"

				code, err := redisdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:    "18512345678",
			wantCode: http.StatusOK,
			wnatBody: ginx.Result{
				Code: 4,
				Msg:  "Send SMS code too fast",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			// Generate a new request and response recorder
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(`{"phone":"`+tc.phone+`"}`)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recoder := httptest.NewRecorder()
			server.ServeHTTP(recoder, req)

			// Verify the response
			assert.Equal(t, tc.wantCode, recoder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res ginx.Result
			err = json.NewDecoder(recoder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wnatBody, res)
		})
	}
}
