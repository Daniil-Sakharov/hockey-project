package infrastructure

import (
	"github.com/robfig/cron/v3"
)

// CronAdapter обёртка над robfig/cron
type CronAdapter struct {
	cron *cron.Cron
}

// NewCronAdapter создаёт новый адаптер
func NewCronAdapter() *CronAdapter {
	return &CronAdapter{
		cron: cron.New(),
	}
}

// AddJob добавляет задачу по cron-выражению
func (a *CronAdapter) AddJob(cronExpr string, handler func()) (cron.EntryID, error) {
	return a.cron.AddFunc(cronExpr, handler)
}

// Start запускает планировщик
func (a *CronAdapter) Start() {
	a.cron.Start()
}

// Stop останавливает планировщик и ждёт завершения задач
func (a *CronAdapter) Stop() {
	ctx := a.cron.Stop()
	<-ctx.Done()
}

// Entries возвращает список зарегистрированных задач
func (a *CronAdapter) Entries() []cron.Entry {
	return a.cron.Entries()
}
