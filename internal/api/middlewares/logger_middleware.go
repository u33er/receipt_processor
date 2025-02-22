package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)

			err := next(c)

			duration := time.Since(start)
			res := c.Response()

			log := logger.With(
				zap.String("request_id", requestID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", res.Status),
				zap.Duration("duration", duration),
				zap.Error(err),
			)

			if c.Request().Method != http.MethodGet && c.Request().Method != http.MethodOptions {
				log.Info("HTTP Request")
			}

			return err
		}
	}
}
