package infrastructure

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// SchedulerMetrics метрики для планировщика
type SchedulerMetrics struct {
	jobDuration      metric.Float64Histogram
	jobTotal         metric.Int64Counter
	tournamentsTotal metric.Int64Counter
	recordsTotal     metric.Int64Counter
	errorsTotal      metric.Int64Counter
}

// NewSchedulerMetrics создаёт метрики для scheduler
func NewSchedulerMetrics() (*SchedulerMetrics, error) {
	meter := otel.Meter("hockey-scheduler")

	jobDuration, err := meter.Float64Histogram(
		"scheduler_job_duration_seconds",
		metric.WithDescription("Duration of scheduler job execution"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	jobTotal, err := meter.Int64Counter(
		"scheduler_job_total",
		metric.WithDescription("Total number of scheduler job executions"),
	)
	if err != nil {
		return nil, err
	}

	tournamentsTotal, err := meter.Int64Counter(
		"scheduler_tournaments_parsed_total",
		metric.WithDescription("Total number of tournaments parsed"),
	)
	if err != nil {
		return nil, err
	}

	recordsTotal, err := meter.Int64Counter(
		"scheduler_records_saved_total",
		metric.WithDescription("Total number of records saved"),
	)
	if err != nil {
		return nil, err
	}

	errorsTotal, err := meter.Int64Counter(
		"scheduler_errors_total",
		metric.WithDescription("Total number of errors"),
	)
	if err != nil {
		return nil, err
	}

	return &SchedulerMetrics{
		jobDuration:      jobDuration,
		jobTotal:         jobTotal,
		tournamentsTotal: tournamentsTotal,
		recordsTotal:     recordsTotal,
		errorsTotal:      errorsTotal,
	}, nil
}

// RecordJobExecution записывает метрики выполнения задачи
func (m *SchedulerMetrics) RecordJobExecution(ctx context.Context, jobName string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}

	attrs := []attribute.KeyValue{
		attribute.String("job_name", jobName),
		attribute.String("status", status),
	}

	m.jobDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
	m.jobTotal.Add(ctx, 1, metric.WithAttributes(attrs...))

	if !success {
		m.errorsTotal.Add(ctx, 1, metric.WithAttributes(
			attribute.String("job_name", jobName),
			attribute.String("error_type", "job_failed"),
		))
	}
}

// RecordTournamentsParsed записывает количество спарсенных турниров
func (m *SchedulerMetrics) RecordTournamentsParsed(ctx context.Context, source string, count int64) {
	m.tournamentsTotal.Add(ctx, count, metric.WithAttributes(
		attribute.String("source", source),
	))
}

// RecordRecordsSaved записывает количество сохранённых записей
func (m *SchedulerMetrics) RecordRecordsSaved(ctx context.Context, recordType string, count int64) {
	m.recordsTotal.Add(ctx, count, metric.WithAttributes(
		attribute.String("type", recordType),
	))
}

// RecordError записывает ошибку
func (m *SchedulerMetrics) RecordError(ctx context.Context, jobName, errorType string) {
	m.errorsTotal.Add(ctx, 1, metric.WithAttributes(
		attribute.String("job_name", jobName),
		attribute.String("error_type", errorType),
	))
}
