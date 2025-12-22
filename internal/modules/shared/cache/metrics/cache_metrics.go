package metrics

import (
	"context"
)

// CacheMetrics метрики кеша
type CacheMetrics struct {
	// Используем pkg/metrics для базовых метрик
}

// NewCacheMetrics создает метрики кеша
func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{}
}

// RecordHit записывает попадание в кеш
func (m *CacheMetrics) RecordHit(ctx context.Context, tier string) {
	// Интеграция с pkg/metrics
}

// RecordMiss записывает промах кеша
func (m *CacheMetrics) RecordMiss(ctx context.Context, tier string) {
	// Интеграция с pkg/metrics
}
