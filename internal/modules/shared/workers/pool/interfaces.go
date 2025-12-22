package pool

import (
	"context"
	"time"
)

// Task интерфейс задачи для выполнения
type Task interface {
	Execute(ctx context.Context) Result
	ID() string
	Priority() int // 0 = низкий, 10 = высокий
}

// Result результат выполнения задачи
type Result interface {
	TaskID() string
	Error() error
	Data() interface{}
	Duration() int64 // в миллисекундах
}

// BaseTask базовая реализация Task
type BaseTask struct {
	id       string
	priority int
	handler  func(ctx context.Context) (interface{}, error)
}

// NewBaseTask создает базовую задачу
func NewBaseTask(id string, priority int, handler func(ctx context.Context) (interface{}, error)) *BaseTask {
	return &BaseTask{
		id:       id,
		priority: priority,
		handler:  handler,
	}
}

func (t *BaseTask) ID() string    { return t.id }
func (t *BaseTask) Priority() int { return t.priority }

func (t *BaseTask) Execute(ctx context.Context) Result {
	start := getCurrentTimeMs()
	data, err := t.handler(ctx)
	duration := getCurrentTimeMs() - start

	return &BaseResult{
		taskID:   t.id,
		data:     data,
		err:      err,
		duration: duration,
	}
}

// BaseResult базовая реализация Result
type BaseResult struct {
	taskID   string
	data     interface{}
	err      error
	duration int64
}

func (r *BaseResult) TaskID() string    { return r.taskID }
func (r *BaseResult) Error() error      { return r.err }
func (r *BaseResult) Data() interface{} { return r.data }
func (r *BaseResult) Duration() int64   { return r.duration }

// getCurrentTimeMs возвращает текущее время в миллисекундах
func getCurrentTimeMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
