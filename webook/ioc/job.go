package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
	"webook/webook/internal/job"
	"webook/webook/internal/service"
	"webook/webook/pkg/logger"
)

func InitRankingJob(svc service.RankingService,
	lockClient *rlock.Client, l logger.Logger) *job.RankingJob {
	return job.NewRankingJob(svc, 30*time.Second, lockClient, l)
}

func InitJobs(l logger.Logger, rjob *job.RankingJob) *cron.Cron {
	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "webook",
		Subsystem: "cronjob",
		Name:      "ranking_cron_job",
		Objectives: map[float64]float64{
			0.5:   0.05,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1m", builder.Build(rjob))
	if err != nil {
		panic(err)
	}
	return expr
}
