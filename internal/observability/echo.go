package observability

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TraceMiddleware instruments each request with an OpenTelemetry span and
// propagates the enriched context downstream.
func TraceMiddleware(serviceName string) echo.MiddlewareFunc {
	tracer := otel.Tracer("http.server")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx, span := tracer.Start(req.Context(), req.Method+" "+c.Path())
			span.SetAttributes(
				attribute.String("http.method", req.Method),
				attribute.String("http.route", c.Path()),
				attribute.String("http.scheme", c.Scheme()),
				attribute.String("http.target", req.URL.Path),
				attribute.String("service.name", serviceName),
			)
			defer span.End()

			c.SetRequest(req.WithContext(ctx))

			if err := next(c); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				slog.Error("request failed", slog.String("path", c.Path()), slog.Any("err", err))
				return err
			}

			span.SetAttributes(attribute.Int("http.status_code", c.Response().Status))
			if c.Response().Status >= 500 {
				span.SetStatus(codes.Error, "server error")
			} else {
				span.SetStatus(codes.Ok, "ok")
			}

			return nil
		}
	}
}
