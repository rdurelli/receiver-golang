package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_mp4_total",
		Help: "The total number of processed mp4 files.",
	})
)

type PrometheusService struct {
}

func (prometheusService PrometheusService) Increment() {
	opsProcessed.Inc()
}
