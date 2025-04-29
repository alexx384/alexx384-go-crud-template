package middleware

import (
	"crud/internal/util/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	slogctx "github.com/veqryn/slog-context"
	"log/slog"
	"time"
)

// JSONLogMiddleware logs a gin HTTP request in JSON format, with some additional custom key/values
func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		ctx := slogctx.Append(c.Request.Context(),
			slog.String("request_id", uuid.NewString()),
			slog.String("client_ip", request.GetClientIP(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("referer", c.Request.Referer()))

		c.Request = c.Request.WithContext(ctx)

		// Process Request
		c.Next()

		// Stop timer
		duration := request.GetDurationInMilliseconds(start)

		logger := slog.With(
			slog.Float64("duration", duration),
			slog.Int("status", c.Writer.Status()))

		if c.Writer.Status() >= 500 {
			logger.ErrorContext(ctx, c.Errors.String())
		} else {
			logger.InfoContext(ctx, "")
		}
	}
}
