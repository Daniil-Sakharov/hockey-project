package domain

import "time"

// Priority приоритет парсинга турнира
type Priority string

const (
	PriorityActive  Priority = "ACTIVE"  // is_ended=false или end_date IS NULL
	PriorityRecent  Priority = "RECENT"  // завершён < 1 месяца
	PriorityMedium  Priority = "MEDIUM"  // завершён 1-6 месяцев
	PriorityOld     Priority = "OLD"     // завершён 6-12 месяцев
	PriorityArchive Priority = "ARCHIVE" // завершён > 1 года
)

// Job представляет задачу планировщика
type Job struct {
	Name    string
	Cron    string
	Timeout time.Duration
	Handler func() error
}

// JobResult результат выполнения задачи
type JobResult struct {
	JobName        string
	StartedAt      time.Time
	EndedAt        time.Time
	Success        bool
	Error          error
	ItemsProcessed int
}
