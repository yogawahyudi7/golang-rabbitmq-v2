package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"golang-rabbitmq-v2/pkg/logger"
)

func Logging(log *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.WithContext("gin").WithFields(map[string]interface{}{
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"status_code": param.StatusCode,
			"latency_ms":  param.Latency.Milliseconds(),
			"user_agent":  param.Request.UserAgent(),
			"error":       param.ErrorMessage,
		}).Info("HTTP Request")
		
		return ""
	})
}

func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		fields := map[string]interface{}{
			"client_ip":     clientIP,
			"method":        method,
			"path":          path,
			"status_code":   statusCode,
			"latency_ms":    latency.Milliseconds(),
			"response_size": c.Writer.Size(),
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		logLevel := "info"
		if statusCode >= 400 && statusCode < 500 {
			logLevel = "warn"
		} else if statusCode >= 500 {
			logLevel = "error"
		}

		switch logLevel {
		case "warn":
			log.WithContext("gin").WithFields(fields).Warn("HTTP Request")
		case "error":
			log.WithContext("gin").WithFields(fields).Error("HTTP Request")
		default:
			log.WithContext("gin").WithFields(fields).Info("HTTP Request")
		}
	}
}