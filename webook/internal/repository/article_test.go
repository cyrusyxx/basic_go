package repository

import (
	"go.uber.org/mock/gomock"
	"testing"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/dao"
)

func TestCachedArticleRepository_Sync(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.ArticleAuthorDAO,
			dao.ArticleReaderDAO)

		arti domain.Article

		wantId  int64
		wantErr error
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
