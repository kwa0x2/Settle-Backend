package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	metrics "github.com/kwa0x2/Settle-Backend/monitoring/prometheus"
	"time"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		statusCode := c.Writer.Status()
		metrics.HttpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(duration)
		metrics.HttpRequestCount.WithLabelValues(c.Request.Method, c.Request.URL.Path, fmt.Sprintf("%d", statusCode)).Inc()
	}
}
