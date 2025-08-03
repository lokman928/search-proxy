package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lokman928/search-proxy/internal/common/logger"
)

func LogRequest(l logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := make(map[string]any)
		fields["method"] = c.Request.Method
		fields["path"] = c.Request.URL.Path

		start := time.Now()

		c.Next()

		fields["latency"] = time.Since(start).String()
		fields["status_code"] = c.Writer.Status()

		l.WithFields(fields).Info("Finish handling request")
	}
}
