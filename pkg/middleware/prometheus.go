package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"service", "method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "method", "path", "status"},
	)

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"service", "method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(responseStatus)
}

func PrometheusMiddleware() gin.HandlerFunc {
	service := "auth"
	return func(c *gin.Context) {
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.Request.URL.Path

		requestCount.WithLabelValues(service, method, path, status).Inc()
		responseStatus.WithLabelValues(service, method, path, status).Inc()
		promhttp.InstrumentHandlerDuration(requestDuration.MustCurryWith(prometheus.Labels{"service": service, "method": method, "path": path, "status": status}), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do nothing, InstrumentHandlerDuration already updates requestDuration.
		})).ServeHTTP(c.Writer, c.Request)
	}
}
