package metrics

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	AttrServiceName = "service.name"
	AttrMethod      = "method"
	AttrStatus      = "status"
	AttrStatusCode  = "status_code"
)

type GRPCMetrics struct {
	RequestsTotal     metric.Int64Counter
	RequestDuration   metric.Float64Histogram
	ActiveConnections metric.Int64UpDownCounter
}

var (
	grpcMetrics     *GRPCMetrics
	grpcMetricsOnce sync.Once
)

func InitGRPCMetrics(ctx context.Context) error {
	var err error

	grpcMetricsOnce.Do(func() {
		meter := otel.Meter("platform.grpc")
		grpcMetrics = &GRPCMetrics{}

		grpcMetrics.RequestsTotal, err = meter.Int64Counter(
			"grpc.server.requests.total",
			metric.WithDescription("Total number of gRPC requests"),
			metric.WithUnit("{request}"),
		)
		if err != nil {
			return
		}

		grpcMetrics.RequestDuration, err = meter.Float64Histogram(
			"grpc.server.request.duration",
			metric.WithDescription("Duration of gRPC requests"),
			metric.WithUnit("s"),
			metric.WithExplicitBucketBoundaries(
				0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0,
			),
		)
		if err != nil {
			return
		}

		grpcMetrics.ActiveConnections, err = meter.Int64UpDownCounter(
			"grpc.server.active_connections",
			metric.WithDescription("Number of active gRPC connections"),
			metric.WithUnit("{connection}"),
		)
	})

	return err
}

// GetGRPCMetrics возвращает singleton gRPC метрик
func GetGRPCMetrics() *GRPCMetrics {
	return grpcMetrics
}

// RecordRequest записывает метрику запроса
func (m *GRPCMetrics) RecordRequest(ctx context.Context, serviceName, method, status string, statusCode int) {
	if m.RequestsTotal == nil {
		return
	}

	m.RequestsTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String(AttrServiceName, serviceName),
			attribute.String(AttrMethod, method),
			attribute.String(AttrStatus, status),
			attribute.Int(AttrStatusCode, statusCode),
		),
	)
}

// RecordDuration записывает длительность запроса
func (m *GRPCMetrics) RecordDuration(ctx context.Context, serviceName, method string, durationSeconds float64) {
	if m.RequestDuration == nil {
		return
	}

	m.RequestDuration.Record(ctx, durationSeconds,
		metric.WithAttributes(
			attribute.String(AttrServiceName, serviceName),
			attribute.String(AttrMethod, method),
		),
	)
}

// IncActiveConnections увеличивает счетчик активных соединений
func (m *GRPCMetrics) IncActiveConnections(ctx context.Context, serviceName string) {
	if m.ActiveConnections == nil {
		return
	}

	m.ActiveConnections.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String(AttrServiceName, serviceName),
		),
	)
}

// DecActiveConnections уменьшает счетчик активных соединений
func (m *GRPCMetrics) DecActiveConnections(ctx context.Context, serviceName string) {
	if m.ActiveConnections == nil {
		return
	}

	m.ActiveConnections.Add(ctx, -1,
		metric.WithAttributes(
			attribute.String(AttrServiceName, serviceName),
		),
	)
}
