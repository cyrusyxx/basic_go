package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"webook/webook/internal/domain"
	svcmocks "webook/webook/internal/service/mocks"
)

func TestBatchRankingService_TopN(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (InteractiveService, ArticleService)

		wantArtis []domain.Article
		wantErr   error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) (InteractiveService, ArticleService) {
				// New
				artiSvc := svcmocks.NewMockArticleService(ctrl)
				interSvc := svcmocks.NewMockInteractiveService(ctrl)

				// Mock ListPub()
				artiSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), int64(0), int64(2)).
					Return([]domain.Article{
						{Id: 1, Title: "title1"},
						{Id: 2, Title: "title2"},
					}, nil)
				artiSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), int64(2), int64(2)).
					Return([]domain.Article{
						{Id: 3, Title: "title3"},
						{Id: 4, Title: "title4"},
					}, nil)
				artiSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), int64(4), int64(2)).
					Return([]domain.Article{}, nil)

				// Mock GetByIds()
				interSvc.EXPECT().GetByIds(gomock.Any(), gomock.Any(), []int64{1, 2}).
					Return(map[int64]domain.InteractiveCount{
						1: {LikeCnt: 1},
						2: {LikeCnt: 2},
					}, nil)
				interSvc.EXPECT().GetByIds(gomock.Any(), gomock.Any(), []int64{3, 4}).
					Return(map[int64]domain.InteractiveCount{
						3: {LikeCnt: 3},
						4: {LikeCnt: 4},
					}, nil)
				interSvc.EXPECT().GetByIds(gomock.Any(), gomock.Any(), []int64{}).
					Return(map[int64]domain.InteractiveCount{}, nil)

				return interSvc, artiSvc
			},
			wantErr: nil,
			wantArtis: []domain.Article{
				{Id: 4, Title: "title4"},
				{Id: 3, Title: "title3"},
				{Id: 2, Title: "title2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interSvc, artiSvc := tc.mock(ctrl)
			svc := &BatchRankingService{
				interSvc:  interSvc,
				artiSvc:   artiSvc,
				batchSize: 2,
				scoreFunc: func(likeCnt int64, utime time.Time) float64 {
					return float64(likeCnt)
				},
			}

			// Test
			artis, err := svc.topN(context.Background(), 3)
			assert.Equal(t, tc.wantArtis, artis)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
