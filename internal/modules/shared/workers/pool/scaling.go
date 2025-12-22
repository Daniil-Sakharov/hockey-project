package pool

import "time"

// adaptiveScaling адаптивное масштабирование пула
func (p *ConfigurablePool) adaptiveScaling() {
	ticker := time.NewTicker(p.config.ScaleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.checkAndScale()
		case <-p.ctx.Done():
			return
		}
	}
}

// checkAndScale проверяет загрузку и масштабирует пул
func (p *ConfigurablePool) checkAndScale() {
	queueLength := len(p.tasks)
	bufferSize := cap(p.tasks)
	currentWorkers := len(p.workers)

	// Вычисляем утилизацию очереди
	utilization := float64(queueLength) / float64(bufferSize)

	// Записываем метрику утилизации
	if p.metrics.PoolUtilization != nil {
		p.metrics.PoolUtilization.Record(p.ctx, utilization)
	}

	// Масштабирование вверх
	if utilization > p.config.ScaleThreshold && currentWorkers < p.maxWorkers {
		newWorkers := min(2, p.maxWorkers-currentWorkers) // Добавляем максимум 2 воркера
		p.startWorkers(newWorkers)
	}

	// Масштабирование вниз (если очередь пуста долгое время)
	// TODO: implement graceful worker shutdown when utilization < 0.1
	_ = utilization // suppress unused warning until implemented
}

// min возвращает минимальное значение
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
