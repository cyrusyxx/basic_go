package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callbacks struct {
	vector *prometheus.SummaryVec
}

func NewCallbacks(opts prometheus.SummaryOpts) *Callbacks {
	v := prometheus.NewSummaryVec(opts, []string{"type", "table"})
	prometheus.MustRegister(v)
	return &Callbacks{
		vector: v,
	}
}

func (c *Callbacks) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	}
}

func (c *Callbacks) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tmp, _ := db.Get("start_time")
		start, ok := tmp.(time.Time)
		if ok {
			dur := time.Since(start).Milliseconds()
			c.vector.WithLabelValues(typ, db.Statement.Table).
				Observe(float64(dur))
		}
	}
}
