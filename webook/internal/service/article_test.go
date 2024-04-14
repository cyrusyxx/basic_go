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

func Test_artileService_Publish(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.ArticleRepository

		arti domain.Article

		wantId  int64
		wantErr error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository,
				repository.ArticleReaderRepository) {
				repo := repomocks.NewMockArticleRepository(ctrl)
				return repo
			},
			arti: domain.Article{
				Title:   "My title",
				Content: "My content",
				Author: domain.Author{
					Id:   123,
					Name: "",
				},
			},
			wantId:  123,
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewImplArticleService(tc.mock(ctrl))
			id, err := svc.Publish(context.Background(), tc.arti)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
