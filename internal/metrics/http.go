package metrics

import "github.com/arthureichelberger/goboiler/pkg/prom"

var (
	httpRequestCounter   = prom.CounterVecFactory("http_requests_total", "Number of requests by method, path and status code", []string{"method", "path", "status"})
	httpRequestHistogram = prom.HistogramVecFactory("http_request_duration_milliseconds", "Duration of requests by method, path and status code", []string{"method", "path", "status"})
)

func CountHTTPRequest(method string, path string, status string) {
	httpRequestCounter.WithLabelValues(method, path, status).Inc()
}

func ObserveHTTPRequestDuration(method string, path string, status string, duration float64) {
	httpRequestHistogram.WithLabelValues(method, path, status).Observe(duration)
}
