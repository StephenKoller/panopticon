package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	taskOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_operations_total",
			Help: "Total number of task operations",
		},
		[]string{"operation"},
	)

	activeTasksGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_tasks",
			Help: "Current number of active (incomplete) tasks",
		},
	)

	completedTasksGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "completed_tasks",
			Help: "Current number of completed tasks",
		},
	)
)
