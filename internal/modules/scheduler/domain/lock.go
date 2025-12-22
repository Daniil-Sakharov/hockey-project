package domain

import "time"

// Lock представляет блокировку задачи
type Lock struct {
	JobName     string
	LockedAt    time.Time
	LockedUntil time.Time
	InstanceID  string
}

// IsExpired проверяет истёк ли срок блокировки
func (l *Lock) IsExpired() bool {
	return time.Now().After(l.LockedUntil)
}
