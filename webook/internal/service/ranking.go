package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"math"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
)

type RankingService interface {
	TopN(ctx context.Context, n int64) error

	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	interSvc InteractiveService
	artiSvc  ArticleService

	repo repository.RankingRepository

	batchSize int64

	scoreFunc func(likeCnt int64, utime time.Time) float64
}

func NewBatchRankingService(repo repository.RankingRepository, interSvc InteractiveService,
	artiSvc ArticleService) RankingService {
	return &BatchRankingService{
		interSvc:  interSvc,
		artiSvc:   artiSvc,
		repo:      repo,
		batchSize: 100,
		scoreFunc: func(likeCnt int64, utime time.Time) float64 {
			dur := time.Since(utime).Seconds()
			return float64(likeCnt-1) / math.Pow(dur+2, 1.5)
		},
	}
}

func (s *BatchRankingService) TopN(ctx context.Context, n int64) error {
	artis, err := s.topN(ctx, n)
	if err != nil {
		return err
	}

	// Save in cache
	return s.repo.ReplaceTopN(ctx, artis)
}

func (s *BatchRankingService) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return s.repo.GetTopN(ctx)
}

func (s *BatchRankingService) topN(ctx context.Context, n int64) ([]domain.Article, error) {
	offset := int64(0)
	start := time.Now()
	ddl := start.Add(-7 * 24 * time.Hour)

	// New Priority Queue
	type Score struct {
		score float64
		arti  domain.Article
	}
	pq := queue.NewPriorityQueue[Score](int(n), func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score < dst.score {
			return -1
		} else {
			return 0
		}
	})

	for {
		// Get a Batch Article
		artis, err := s.artiSvc.ListPub(ctx, start, offset, s.batchSize)
		if err != nil {
			return nil, err
		}
		// Get Interactive Count From Batch
		ids := getids(artis)
		interMap, err := s.interSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		// Calculate Score And Enqueue
		for _, arti := range artis {
			inter := interMap[arti.Id]
			score := s.scoreFunc(inter.LikeCnt, arti.Utime)
			item := Score{
				score: score,
				arti:  arti,
			}
			err = pq.Enqueue(item)
			if err == queue.ErrOutOfCapacity {
				// Replace the smallest score
				minItem, _ := pq.Dequeue()
				if minItem.score < score {
					_ = pq.Enqueue(item)
				} else {
					_ = pq.Enqueue(minItem)
				}
			}
		}
		artisLen := int64(len(artis))
		offset = offset + artisLen
		if artisLen < s.batchSize || artis[artisLen-1].Utime.Before(ddl) {
			break
		}
	}

	res := make([]domain.Article, pq.Len())
	for i := pq.Len() - 1; i >= 0; i-- {
		item, _ := pq.Dequeue()
		res[i] = item.arti
	}
	return res, nil
}

// getids get ids by articles
func getids(artis []domain.Article) []int64 {
	if len(artis) == 0 {
		return []int64{}
	}
	var ids []int64
	for _, arti := range artis {
		ids = append(ids, arti.Id)
	}
	return ids
}
