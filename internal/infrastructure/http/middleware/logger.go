// Package middleware contains cross-cutting HTTP concerns (logging, etc.)
// implemented as Gin middleware functions.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RequestLogger returns a Gin middleware that logs every incoming HTTP
// request using structured logging (logrus), including method, path,
// status code, client IP and latency.
func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if raw != "" {
			path = path + "?" + raw
		}

		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
		}).Info("http request handled")
	}
}

