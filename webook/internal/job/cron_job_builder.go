package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
	"webook/webook/pkg/logger"
)

type CronJobBuilder struct {
	l   logger.Logger
	vec *prometheus.SummaryVec
}

func NewCronJobBuilder(l logger.Logger, opt prometheus.SummaryOpts) *CronJobBuilder {
	return &CronJobBuilder{
		l:   l,
		vec: prometheus.NewSummaryVec(opt, []string{"name", "success"}),
	}
}

func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cron.FuncJob(func() {
		start := time.Now()
		err := job.Run()
		if err != nil {
			b.l.Error("failed to run job",
				logger.Error(err),
				logger.String("name", name),
			)
		}
		dur := time.Since(start)
		b.vec.WithLabelValues(name, strconv.FormatBool(err == nil)).Observe(float64(dur.Milliseconds()))
	})
}
