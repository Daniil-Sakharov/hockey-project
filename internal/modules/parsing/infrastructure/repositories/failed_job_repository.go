package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// FailedJob представляет неудачную задачу парсинга
type FailedJob struct {
	ID           int            `db:"id"`
	JobType      string         `db:"job_type"`
	Source       string         `db:"source"`
	ExternalID   string         `db:"external_id"`
	URL          sql.NullString `db:"url"`
	ErrorMessage sql.NullString `db:"error_message"`
	RetryCount   int            `db:"retry_count"`
	MaxRetries   int            `db:"max_retries"`
	NextRetryAt  time.Time      `db:"next_retry_at"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

// FailedJobRepository репозиторий для работы с неудачными задачами
type FailedJobRepository struct {
	db *sqlx.DB
}

// NewFailedJobRepository создаёт новый репозиторий
func NewFailedJobRepository(db *sqlx.DB) *FailedJobRepository {
	return &FailedJobRepository{db: db}
}

// Create создаёт запись о неудачной задаче
func (r *FailedJobRepository) Create(ctx context.Context, job *FailedJob) error {
	query := `
		INSERT INTO failed_parsing_jobs (job_type, source, external_id, url, error_message, retry_count, max_retries, next_retry_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query,
		job.JobType, job.Source, job.ExternalID, job.URL, job.ErrorMessage,
		job.RetryCount, job.MaxRetries, job.NextRetryAt,
	)
	return err
}

// GetPendingRetries возвращает задачи готовые к повторной попытке
func (r *FailedJobRepository) GetPendingRetries(ctx context.Context, limit int) ([]*FailedJob, error) {
	query := `
		SELECT * FROM failed_parsing_jobs 
		WHERE next_retry_at <= NOW() AND retry_count < max_retries
		ORDER BY next_retry_at
		LIMIT $1
	`
	var jobs []*FailedJob
	if err := r.db.SelectContext(ctx, &jobs, query, limit); err != nil {
		return nil, err
	}
	return jobs, nil
}

// IncrementRetry увеличивает счётчик попыток и устанавливает следующее время
func (r *FailedJobRepository) IncrementRetry(ctx context.Context, id int, nextRetryAt time.Time, errMsg string) error {
	query := `
		UPDATE failed_parsing_jobs 
		SET retry_count = retry_count + 1, 
		    next_retry_at = $2, 
		    error_message = $3,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, nextRetryAt, errMsg)
	return err
}

// Delete удаляет задачу (при успешном выполнении)
func (r *FailedJobRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM failed_parsing_jobs WHERE id = $1`, id)
	return err
}

// CountPending возвращает количество ожидающих задач
func (r *FailedJobRepository) CountPending(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM failed_parsing_jobs WHERE retry_count < max_retries`)
	return count, err
}

// CountFailed возвращает количество окончательно неудачных задач
func (r *FailedJobRepository) CountFailed(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM failed_parsing_jobs WHERE retry_count >= max_retries`)
	return count, err
}

// SaveError сохраняет ошибку парсинга
func (r *FailedJobRepository) SaveError(ctx context.Context, jobType, source, externalID, url string, err error) error {
	job := &FailedJob{
		JobType:      jobType,
		Source:       source,
		ExternalID:   externalID,
		URL:          sql.NullString{String: url, Valid: url != ""},
		ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
		RetryCount:   0,
		MaxRetries:   3,
		NextRetryAt:  time.Now().Add(5 * time.Minute),
	}
	return r.Create(ctx, job)
}

// GetRetryInterval возвращает интервал до следующей попытки
func GetRetryInterval(retryCount int) time.Duration {
	switch retryCount {
	case 0:
		return 5 * time.Minute
	case 1:
		return 30 * time.Minute
	case 2:
		return 2 * time.Hour
	default:
		return 24 * time.Hour
	}
}

// CleanupOld удаляет старые записи
func (r *FailedJobRepository) CleanupOld(ctx context.Context, olderThan time.Duration) (int64, error) {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM failed_parsing_jobs WHERE created_at < $1`,
		time.Now().Add(-olderThan),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// GetBySource возвращает задачи по источнику
func (r *FailedJobRepository) GetBySource(ctx context.Context, source string) ([]*FailedJob, error) {
	var jobs []*FailedJob
	err := r.db.SelectContext(ctx, &jobs,
		`SELECT * FROM failed_parsing_jobs WHERE source = $1 ORDER BY created_at DESC`,
		source,
	)
	return jobs, err
}

// MarkAsFailed помечает задачу как окончательно неудачную
func (r *FailedJobRepository) MarkAsFailed(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE failed_parsing_jobs SET retry_count = max_retries, updated_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}

// Exists проверяет существует ли уже такая задача
func (r *FailedJobRepository) Exists(ctx context.Context, jobType, source, externalID string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM failed_parsing_jobs WHERE job_type = $1 AND source = $2 AND external_id = $3 AND retry_count < max_retries`,
		jobType, source, externalID,
	)
	return count > 0, err
}
