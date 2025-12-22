package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/scheduler/domain"
)

// LockRepository репозиторий для работы с блокировками
type LockRepository struct {
	db *sql.DB
}

// NewLockRepository создаёт новый репозиторий блокировок
func NewLockRepository(db *sql.DB) *LockRepository {
	return &LockRepository{db: db}
}

// TryAcquire пытается получить блокировку
func (r *LockRepository) TryAcquire(ctx context.Context, jobName string, timeout time.Duration, instanceID string) (bool, error) {
	query := `
		INSERT INTO scheduler_locks (job_name, locked_at, locked_until, instance_id)
		VALUES ($1, NOW(), NOW() + $2::interval, $3)
		ON CONFLICT (job_name) DO UPDATE 
		SET locked_at = NOW(), 
		    locked_until = NOW() + $2::interval, 
		    instance_id = $3
		WHERE scheduler_locks.locked_until < NOW()
	`

	result, err := r.db.ExecContext(ctx, query, jobName, timeout.String(), instanceID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// Release освобождает блокировку
func (r *LockRepository) Release(ctx context.Context, jobName, instanceID string) error {
	query := `DELETE FROM scheduler_locks WHERE job_name = $1 AND instance_id = $2`
	_, err := r.db.ExecContext(ctx, query, jobName, instanceID)
	return err
}

// Get получает информацию о блокировке
func (r *LockRepository) Get(ctx context.Context, jobName string) (*domain.Lock, error) {
	query := `SELECT job_name, locked_at, locked_until, instance_id FROM scheduler_locks WHERE job_name = $1`

	var lock domain.Lock
	err := r.db.QueryRowContext(ctx, query, jobName).Scan(
		&lock.JobName,
		&lock.LockedAt,
		&lock.LockedUntil,
		&lock.InstanceID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &lock, nil
}

// ReleaseAll освобождает все блокировки инстанса
func (r *LockRepository) ReleaseAll(ctx context.Context, instanceID string) error {
	query := `DELETE FROM scheduler_locks WHERE instance_id = $1`
	_, err := r.db.ExecContext(ctx, query, instanceID)
	return err
}
