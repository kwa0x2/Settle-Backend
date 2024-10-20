package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Histogram of HTTP request durations.",
		},
		[]string{"method", "path"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(HttpRequestCount)
	prometheus.MustRegister(HttpRequestDuration)
}
