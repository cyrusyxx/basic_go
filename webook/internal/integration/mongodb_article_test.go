package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webook/webook/internal/integration/startup"
	"webook/webook/internal/repository/dao"
	ijwt "webook/webook/internal/web/jwt"
)

func TestMongoArticleHandler_Edit(t *testing.T) {
	// Generate context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Init MongoDB And Collection
	mongodb := startup.InitMongoDB()
	col := mongodb.Collection("articles")
	livecol := mongodb.Collection("published_articles")
	err := dao.InitCollection(mongodb)
	assert.NoError(t, err)

	// Init Handler
	node, err := snowflake.NewNode(1)
	assert.NoError(t, err)
	mDao := dao.NewMongoDBArticleDAO(mongodb, node)
	hdl := startup.InitArticleHandler(mDao)

	// Start Server
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("userclaim", ijwt.UserClaims{
			Uid:       123,
			UserAgent: "",
			Ssid:      "",
		})
	})
	hdl.RegisterRoutes(server)

	// Define Testcases
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
				// Find And Decode on &arti
				var arti dao.Article
				err := col.FindOne(ctx, bson.D{{"author_id", 123}}).
					Decode(&arti)

				// Verify
				assert.NoError(t, err)
				assert.True(t, arti.Ctime > 0)
				assert.True(t, arti.Utime > 0)
				assert.True(t, arti.Id != 0)
				assert.Equal(t, "My Title", arti.Title)
				assert.Equal(t, "My Content", arti.Content)
				assert.Equal(t, int64(123), arti.AuthorId)

				// Clear All Test Data
				_, err = col.DeleteMany(ctx, bson.D{})
				assert.NoError(t, err)
				_, err = livecol.DeleteMany(ctx, bson.D{})
				assert.NoError(t, err)
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
				// Just Insert One
				_, err := col.InsertOne(ctx, dao.Article{
					Id:       2,
					Title:    "My Title",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    456,
					Utime:    789,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// Find
				var arti dao.Article
				err := col.FindOne(ctx, bson.D{{"id", 2}}).Decode(&arti)

				// Verify
				assert.NoError(t, err)
				assert.True(t, arti.Ctime == 456)
				assert.True(t, arti.Utime > 789)
				assert.True(t, arti.Id == 2)
				assert.Equal(t, "New Title", arti.Title)
				assert.Equal(t, "New Content", arti.Content)
				assert.Equal(t, int64(123), arti.AuthorId)

				_, err = col.DeleteMany(ctx, bson.D{})
				assert.NoError(t, err)
				_, err = livecol.DeleteMany(ctx, bson.D{})
				assert.NoError(t, err)
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
				// Just Insert One
				_, err := col.InsertOne(ctx, dao.Article{
					Id:       3,
					Title:    "My Title",
					Content:  "My Content",
					AuthorId: 234,
					Ctime:    456,
					Utime:    789,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// Find
				var arti dao.Article
				err := col.FindOne(ctx, bson.D{{"id", 3}}).Decode(&arti)

				// Verify
				assert.NoError(t, err)
				assert.True(t, arti.Ctime == 456)
				assert.True(t, arti.Utime == 789)
				assert.True(t, arti.Id == 3)
				assert.Equal(t, "My Title", arti.Title)
				assert.Equal(t, "My Content", arti.Content)
				assert.Equal(t, int64(234), arti.AuthorId)

				_, err = col.DeleteMany(ctx, bson.D{})
				assert.NoError(t, err)
				_, err = livecol.DeleteMany(ctx, bson.D{})
				assert.NoError(t, err)
			},
			arti: Article{
				Id:      3,
				Title:   "New Title",
				Content: "New Content",
			},
			wantCode: http.StatusOK,
			wantRes:  Result[int64]{},
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
			if res.Data != 0 {
				assert.True(t, res.Data != 0)
			}
		})
	}
}
