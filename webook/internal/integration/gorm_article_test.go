package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/webook/internal/integration/startup"
	"webook/webook/internal/repository/dao"
	ijwt "webook/webook/internal/web/jwt"
)

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func TestArticleHandler_Edit(t *testing.T) {
	mysqldb := startup.InitMysql()
	hdl := startup.InitArticleHandler(dao.NewGORMArticleDAO(mysqldb))
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("userclaim", ijwt.UserClaims{
			Uid:       123,
			UserAgent: "",
			Ssid:      "",
		})
	})
	hdl.RegisterRoutes(server)

	testcases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		arti Article

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "New article",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.Article
				err := mysqldb.Where("id=?", 1).
					First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Id > 0)
				assert.Equal(t, "My Title", art.Title)
				assert.Equal(t, "My Content", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)

				mysqldb.Exec("truncate table `articles`")
			},
			arti: Article{
				Title:   "My Title",
				Content: "My Content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
			},
		},
		{
			name: "Edit article",
			before: func(t *testing.T) {
				err := mysqldb.Create(dao.Article{
					Id:       2,
					Title:    "My Title",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var arti dao.Article
				err := mysqldb.Where("id=?", 2).
					First(&arti).Error
				assert.NoError(t, err)
				assert.True(t, arti.Ctime == 456)
				assert.True(t, arti.Utime > 789)
				assert.True(t, arti.Id == 2)
				assert.Equal(t, "New Title", arti.Title)
				assert.Equal(t, "New Content", arti.Content)
				assert.Equal(t, int64(123), arti.AuthorId)

				mysqldb.Exec("truncate table `articles`")
			},
			arti: Article{
				Id:      2,
				Title:   "New Title",
				Content: "New Content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
			},
		},
		{
			name: "Edit others article",
			before: func(t *testing.T) {
				err := mysqldb.Create(dao.Article{
					Id:       3,
					Title:    "My Title",
					Content:  "My Content",
					AuthorId: 234,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var arti dao.Article
				err := mysqldb.Where("id=?", 3).
					First(&arti).Error
				assert.NoError(t, err)
				assert.True(t, arti.Ctime == 456)
				assert.True(t, arti.Utime == 789)
				assert.True(t, arti.Id == 3)
				assert.Equal(t, "My Title", arti.Title)
				assert.Equal(t, "My Content", arti.Content)
				assert.Equal(t, int64(234), arti.AuthorId)

				mysqldb.Exec("truncate table `articles`")
			},
			arti: Article{
				Id:      3,
				Title:   "My Title",
				Content: "My Content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg: "System Error",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			// Generate a new request and response recorder
			reBody, err := json.Marshal(tc.arti)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/article/edit",
				bytes.NewReader(reBody))

			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			// Verify the response
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
