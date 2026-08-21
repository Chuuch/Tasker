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
		[]string{"method", "route", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "tasker_http_request_duration_seconds",
			Help: "HTTP request duration in seconds.",
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tasker_http_requests_in_flight",
			Help: "Number of HTTP requests currenty being processed by Tasker.",
		},
	)
)

func Register() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		HTTPRequestsInFlight,
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

func RequestStarted() {
	HTTPRequestsInFlight.Inc()
}

func RequestFinished() {
	HTTPRequestsInFlight.Dec()
}
