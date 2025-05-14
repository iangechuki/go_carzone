package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)


var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of all HTTP requests",
		},
		[]string{"path", "method"},
	)
	statusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Count of HTTP requests by status code",
		},
		[]string{"path", "method", "status_code"},
	)
)
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}
func init(){
	prometheus.MustRegister(requestCounter, requestDuration, statusCounter)
}
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(ww, r)
		// duration of req
		duration := time.Since(start).Seconds()
		requestCounter.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(ww.statusCode)).Inc()
		requestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
		statusCounter.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(ww.statusCode)).Inc()
		
	})
}
func (rw *responseWriter)WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}