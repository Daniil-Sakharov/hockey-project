package retry

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

type JobType string

const (
	JobTypeTournament JobType = "tournament"
	JobTypeTeam       JobType = "team"
	JobTypePlayer     JobType = "player"
)

type FailedJob struct {
	ID           int       `db:"id"`
	JobType      JobType   `db:"job_type"`
	Source       string    `db:"source"`
	ExternalID   string    `db:"external_id"`
	URL          string    `db:"url"`
	ErrorMessage string    `db:"error_message"`
	RetryCount   int       `db:"retry_count"`
	MaxRetries   int       `db:"max_retries"`
	NextRetryAt  time.Time `db:"next_retry_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Manager struct {
	db         *sqlx.DB
	maxRetries int
	baseDelay  time.Duration
}

func NewManager(db *sqlx.DB, maxRetries int, baseDelay time.Duration) *Manager {
	return &Manager{
		db:         db,
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
	}
}

// AddFailedJob добавляет неудачную задачу в очередь retry
func (m *Manager) AddFailedJob(ctx context.Context, jobType JobType, source, externalID, url string, err error) error {
	nextRetry := time.Now().Add(m.baseDelay)
	
	query := `
		INSERT INTO failed_parsing_jobs (job_type, source, external_id, url, error_message, max_retries, next_retry_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, dbErr := m.db.ExecContext(ctx, query, jobType, source, externalID, url, err.Error(), m.maxRetries, nextRetry)
	if dbErr != nil {
		logger.Error(ctx, "Failed to add retry job", zap.Error(dbErr))
		return dbErr
	}
	
	logger.Warn(ctx, "Added failed job for retry", 
		zap.String("type", string(jobType)),
		zap.String("source", source),
		zap.String("external_id", externalID),
		zap.Error(err))
	
	return nil
}

// GetJobsForRetry получает задачи готовые для повторной попытки
func (m *Manager) GetJobsForRetry(ctx context.Context, source string, limit int) ([]FailedJob, error) {
	query := `
		SELECT id, job_type, source, external_id, url, error_message, retry_count, max_retries, next_retry_at, created_at, updated_at
		FROM failed_parsing_jobs 
		WHERE source = $1 AND next_retry_at <= NOW() AND retry_count < max_retries
		ORDER BY next_retry_at ASC
		LIMIT $2`
	
	var jobs []FailedJob
	err := m.db.SelectContext(ctx, &jobs, query, source, limit)
	return jobs, err
}

// MarkJobRetried обновляет счетчик попыток и время следующей попытки
func (m *Manager) MarkJobRetried(ctx context.Context, jobID int, success bool, err error) error {
	if success {
		// Удаляем успешно обработанную задачу
		_, dbErr := m.db.ExecContext(ctx, "DELETE FROM failed_parsing_jobs WHERE id = $1", jobID)
		return dbErr
	}
	
	// Увеличиваем счетчик и устанавливаем следующую попытку
	nextRetry := time.Now().Add(m.calculateDelay(1)) // Exponential backoff можно добавить позже
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	
	query := `
		UPDATE failed_parsing_jobs 
		SET retry_count = retry_count + 1, 
		    next_retry_at = $2,
		    error_message = $3,
		    updated_at = NOW()
		WHERE id = $1`
	
	_, dbErr := m.db.ExecContext(ctx, query, jobID, nextRetry, errorMsg)
	return dbErr
}

// CleanupOldJobs удаляет старые неудачные задачи
func (m *Manager) CleanupOldJobs(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	
	result, err := m.db.ExecContext(ctx, 
		"DELETE FROM failed_parsing_jobs WHERE created_at < $1 OR retry_count >= max_retries", 
		cutoff)
	
	if err == nil {
		if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
			logger.Info(ctx, "Cleaned up old failed jobs", zap.Int64("count", rowsAffected))
		}
	}
	
	return err
}

func (m *Manager) calculateDelay(retryCount int) time.Duration {
	// Простая логика: базовая задержка * retry_count
	// Можно улучшить до exponential backoff
	return m.baseDelay * time.Duration(retryCount+1)
}
