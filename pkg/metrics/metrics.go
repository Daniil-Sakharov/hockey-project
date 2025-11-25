package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTimeout = 60 * time.Second // Увеличено для OTEL Collector remote write
)

var (
	otlpExporter       *otlpmetricgrpc.Exporter
	prometheusExporter *prometheus.Exporter
	meterProvider      *metric.MeterProvider
)

type Config interface {
	ServiceName() string
	ServiceVersion() string
	Environment() string
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}

func InitProvider(ctx context.Context, cfg Config) error {
	var err error

	otlpExporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()), // в проде с TLS
		otlpmetricgrpc.WithTimeout(defaultTimeout),                   // Таймаут 60с для экспорта метрик
	)
	if err != nil {
		return errors.Wrap(err, "failed to create metrics exporter")
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName()),
			attribute.String("service.version", cfg.ServiceVersion()),
			attribute.String("deployment.environment", cfg.Environment()),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to set metrics attributes")
	}

	meterProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				otlpExporter,
				metric.WithInterval(cfg.CollectorInterval()),
				metric.WithTimeout(defaultTimeout), // Таймаут для reader
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	return nil
}

// InitProviderWithPrometheus инициализирует метрики с OTLP и Prometheus exporters
// Prometheus exporter используется для scraping метрик напрямую из сервиса
func InitProviderWithPrometheus(ctx context.Context, cfg Config) error {
	var err error

	// Создаем Prometheus exporter
	prometheusExporter, err = prometheus.New()
	if err != nil {
		return errors.Wrap(err, "failed to create prometheus exporter")
	}

	// Создаем OTLP exporter
	otlpExporter, err = otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()),
		otlpmetricgrpc.WithTimeout(defaultTimeout), // Таймаут 60с для экспорта метрик
	)
	if err != nil {
		return errors.Wrap(err, "failed to create otlp exporter")
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName()),
			attribute.String("service.version", cfg.ServiceVersion()),
			attribute.String("deployment.environment", cfg.Environment()),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to set metrics attributes")
	}

	// Создаем MeterProvider с обоими readers
	meterProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(prometheusExporter), // Для Prometheus scraping
		metric.WithReader(
			metric.NewPeriodicReader(
				otlpExporter,
				metric.WithInterval(cfg.CollectorInterval()),
				metric.WithTimeout(defaultTimeout), // Таймаут для reader
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	return nil
}

// MetricsHandler возвращает HTTP handler для Prometheus метрик
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func GetMeterProvider() *metric.MeterProvider {
	return meterProvider
}

func Shutdown(ctx context.Context) error {
	if meterProvider == nil && otlpExporter == nil && prometheusExporter == nil {
		return nil
	}

	var err error

	if meterProvider != nil {
		err = meterProvider.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown meter provider")
		}
	}

	if otlpExporter != nil {
		err = otlpExporter.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown otlp exporter")
		}
	}

	if prometheusExporter != nil {
		err = prometheusExporter.Shutdown(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to shutdown prometheus exporter")
		}
	}

	return nil
}
