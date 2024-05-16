package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Builder struct {
	// WARN: Only "_" and "Aa" is valid
	Namespace string
	Subsystem string
	Name      string

	InstanceId string
}

func NewBuilder(namespace, subsystem, name, instanceId string) *Builder {
	return &Builder{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       name,
		InstanceId: instanceId,
	}
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_response_time",

		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)
	prometheus.MustRegister(vector)
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := c.Request.Method
			pattern := c.FullPath()
			status := c.Writer.Status()
			vector.WithLabelValues(method, pattern, strconv.Itoa(status)).
				Observe(float64(duration))
		}()
		c.Next()
	}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_active_request",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)
	return func(c *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		c.Next()
	}
}
