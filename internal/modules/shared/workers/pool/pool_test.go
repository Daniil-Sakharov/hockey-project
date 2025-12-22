package pool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestConfigurablePool_Basic(t *testing.T) {
	ctx := context.Background()

	config := Config{
		Name:        "test-pool",
		WorkerCount: 2,
		MaxWorkers:  4,
		BufferSize:  10,
		TaskTimeout: 5 * time.Second,
	}

	pool, err := NewConfigurablePool(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	pool.Start()

	// Создаем простую задачу
	task := NewBaseTask("test-1", 5, func(ctx context.Context) (interface{}, error) {
		return "success", nil
	})

	// Отправляем задачу
	pool.Submit(task)

	// Получаем результат
	select {
	case result := <-pool.Results():
		if result.Error() != nil {
			t.Errorf("Task failed: %v", result.Error())
		}
		if result.Data() != "success" {
			t.Errorf("Expected 'success', got %v", result.Data())
		}
	case <-time.After(2 * time.Second):
		t.Error("Task timeout")
	}
}

func TestConfigurablePool_MultipleTask(t *testing.T) {
	ctx := context.Background()

	config := Config{
		Name:        "multi-test-pool",
		WorkerCount: 3,
		MaxWorkers:  6,
		BufferSize:  20,
		TaskTimeout: 5 * time.Second,
	}

	pool, err := NewConfigurablePool(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	pool.Start()

	taskCount := 10

	// Отправляем несколько задач
	for i := 0; i < taskCount; i++ {
		task := NewBaseTask(
			fmt.Sprintf("task-%d", i),
			5,
			func(ctx context.Context) (interface{}, error) {
				time.Sleep(100 * time.Millisecond) // Имитируем работу
				return "done", nil
			},
		)
		pool.Submit(task)
	}

	// Собираем результаты
	completed := 0
	timeout := time.After(5 * time.Second)

	for completed < taskCount {
		select {
		case result := <-pool.Results():
			if result.Error() != nil {
				t.Errorf("Task %s failed: %v", result.TaskID(), result.Error())
			}
			completed++
		case <-timeout:
			t.Errorf("Timeout: completed %d/%d tasks", completed, taskCount)
			return
		}
	}
}
