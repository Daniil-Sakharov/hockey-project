package tracing

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	HTTPTraceIDHeader = "X-Trace-ID"
)

func createHTTPSpanAttributes(r *http.Request) []trace.SpanStartOption {
	return []trace.SpanStartOption{
		trace.WithAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.URLPath(r.URL.Path),
			semconv.HostName(r.Host),
			semconv.URLScheme(r.URL.Scheme),
			semconv.UserAgentName(r.UserAgent()),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	}
}

func HTTPHandlerMiddleware(serviceName string) func(http.Handler) http.Handler {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			spanName := fmt.Sprintf("%s.%s", r.Method, r.URL.Path)

			ctx, span := tracer.Start(
				ctx,
				spanName,
				createHTTPSpanAttributes(r)...,
			)
			defer span.End()

			wrw := &traceResponseWriter{
				ResponseWriter: w,
				span:           span,
				headerAdded:    false,
			}

			next.ServeHTTP(wrw, r.WithContext(ctx))

			span.SetAttributes(semconv.OTelStatusCodeKey.Int(wrw.statusCode))
		})
	}
}
