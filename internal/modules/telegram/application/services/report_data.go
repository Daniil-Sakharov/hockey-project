package services

import "context"

// DataCollector собирает все данные для отчёта
type DataCollector struct {
	repo ReportRepository
}

// NewDataCollector создает новый DataCollector
func NewDataCollector(repo ReportRepository) *DataCollector {
	return &DataCollector{repo: repo}
}

// CollectFullReport собирает все данные для полного отчёта игрока
func (dc *DataCollector) CollectFullReport(ctx context.Context, playerID string) (*FullPlayerReport, error) {
	// Используем существующий метод GetFullReport из ReportRepository
	return dc.repo.GetFullReport(ctx, playerID)
}
