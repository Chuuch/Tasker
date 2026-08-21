package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasker_http_requests_total",
			Help: "Total number of HTTP requests processed by Tasker.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "tasker_http_request_duration_seconds",
			Help: "HTTP request duration in seconds.",
		},
		[]string{"method", "path"},
	)
)

func Register() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
	)
}

func ObserveRequest(method, path string, status int, duration time.Duration) {
	HTTPRequestsTotal.WithLabelValues(
		method,
		path,
		strconv.Itoa(status),
	).Inc()

	HTTPRequestDuration.WithLabelValues(
		method,
		path,
	).Observe(duration.Seconds())
}
