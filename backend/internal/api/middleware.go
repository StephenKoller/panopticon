package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("http")
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		ctx, span := tracer.Start(
			c.Request.Context(),
			fmt.Sprintf("HTTP %s %s", c.Request.Method, path),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.path", path),
			),
		)
		defer span.End()

		// Store the span in the context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Add response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		// Record error if any
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors[0])
		}
	}
}
