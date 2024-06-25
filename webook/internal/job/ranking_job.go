package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"sync"
	"time"
	"webook/webook/internal/service"
	"webook/webook/pkg/logger"
)

type RankingJob struct {
	svc        service.RankingService
	timeout    time.Duration
	lockClient *rlock.Client
	key        string

	l logger.Logger

	localLock *sync.Mutex
	lock      *rlock.Lock
}

func NewRankingJob(svc service.RankingService,
	timeout time.Duration, lockClient *rlock.Client, l logger.Logger) *RankingJob {
	return &RankingJob{
		svc:        svc,
		timeout:    timeout,
		lockClient: lockClient,
		key:        "job:ranking",
		l:          l,
		localLock:  &sync.Mutex{},
	}
}

func (j *RankingJob) Name() string {
	return "RankingJob"
}

func (j *RankingJob) Run() error {
	j.localLock.Lock()
	lock := j.lock
	if lock == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		lock, err := j.lockClient.Lock(ctx, j.key, j.timeout,
			&rlock.FixIntervalRetry{
				Interval: 100 * time.Millisecond,
				Max:      3,
			}, time.Second)
		if err != nil {
			j.localLock.Unlock()
			j.l.Warn("failed to get ranking job lock", logger.Error(err))
			return nil
		}
		j.lock = lock
		j.localLock.Unlock()
		go func() {
			er := lock.AutoRefresh(j.timeout/2, j.timeout)
			if er != nil {
				j.localLock.Lock()
				j.lock = nil
				j.localLock.Unlock()
			}
		}()
	}
	ctx, cancel := context.WithTimeout(context.Background(), j.timeout)
	defer cancel()
	return j.svc.TopN(ctx, 100)
}

//func (j *RankingJob) Run() error {
//	// Get Lock
//	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
//	lock, err := j.lockClient.Lock(ctx, j.key, j.timeout,
//		&rlock.FixIntervalRetry{
//			Interval: 100 * time.Millisecond,
//			Max:      3,
//		}, time.Second)
//	cancel()
//	if err != nil {
//		return err
//	}
//	// Release Lock
//	defer func() {
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//		defer cancel()
//		err := lock.Unlock(ctx)
//		if err != nil {
//			j.l.Error("failed to unlock ranking job lock", logger.Error(err))
//		}
//	}()
//
//	// Run
//	ctx, cancel = context.WithTimeout(context.Background(), j.timeout)
//	defer cancel()
//	return j.svc.TopN(ctx, 100)
//}
