package handlers

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"strings"
)

func Logging(logger *slog.Logger) gin.HandlerFunc {
	return sloggin.NewWithConfig(logger, sloggin.Config{
		WithUserAgent:      false,
		WithRequestID:      false,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: true,
		WithSpanID:         false,
		WithTraceID:        false,
		Filters:            nil,
	})
}

func LoggingWithNoHealth(logger *slog.Logger) gin.HandlerFunc {
	return sloggin.NewWithConfig(logger, sloggin.Config{
		WithUserAgent:      false,
		WithRequestID:      false,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: true,
		WithSpanID:         false,
		WithTraceID:        false,
		Filters:            []sloggin.Filter{IgnoreHealthFilter},
	})
}

func IgnoreHealthFilter(c *gin.Context) bool {
	reqUri := strings.ToLower(c.Request.RequestURI)
	if strings.Contains(reqUri, "healthz") {
		return false
	}
	return true
}
