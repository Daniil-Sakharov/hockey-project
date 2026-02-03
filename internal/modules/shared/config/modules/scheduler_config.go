package modules

import (
	"fmt"
	"os"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

// JobWithName структура для хранения job с именем (для сортировки)
type JobWithName struct {
	Name string
	Job  JobConfig
}

// SchedulerConfig конфигурация планировщика
type SchedulerConfig struct {
	BootstrapMode  bool                 `yaml:"bootstrap_mode"`
	RunImmediately bool                 `yaml:"run_immediately"`
	Jobs           map[string]JobConfig `yaml:"jobs"`
}

// JobConfig конфигурация отдельной задачи
type JobConfig struct {
	Cron           string        `yaml:"cron"`
	Enabled        bool          `yaml:"enabled"`
	Timeout        time.Duration `yaml:"timeout"`
	MaxTournaments int           `yaml:"max_tournaments"`
	Order          int           `yaml:"order"` // Порядок выполнения (меньше = раньше)
}

// LoadSchedulerConfig загружает конфигурацию из YAML файла
func LoadSchedulerConfig(path string) (*SchedulerConfig, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is from trusted config
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var wrapper struct {
		Scheduler SchedulerConfig `yaml:"scheduler"`
	}

	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := wrapper.Scheduler.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &wrapper.Scheduler, nil
}

// Validate проверяет валидность конфигурации
func (c *SchedulerConfig) Validate() error {
	if len(c.Jobs) == 0 {
		return fmt.Errorf("no jobs configured")
	}

	for name, job := range c.Jobs {
		if job.Cron == "" {
			return fmt.Errorf("job %s: cron expression required", name)
		}
		if job.Timeout <= 0 {
			return fmt.Errorf("job %s: timeout must be positive", name)
		}
	}

	return nil
}

// GetJob возвращает конфигурацию задачи по имени
func (c *SchedulerConfig) GetJob(name string) (JobConfig, bool) {
	job, ok := c.Jobs[name]
	return job, ok
}

// EnabledJobs возвращает список включённых задач (deprecated: используй EnabledJobsOrdered)
func (c *SchedulerConfig) EnabledJobs() map[string]JobConfig {
	result := make(map[string]JobConfig)
	for name, job := range c.Jobs {
		if job.Enabled {
			result[name] = job
		}
	}
	return result
}

// EnabledJobsOrdered возвращает список включённых задач, отсортированный по полю Order
func (c *SchedulerConfig) EnabledJobsOrdered() []JobWithName {
	var jobs []JobWithName
	for name, job := range c.Jobs {
		if job.Enabled {
			jobs = append(jobs, JobWithName{Name: name, Job: job})
		}
	}

	// Сортируем по Order (меньше = раньше)
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Job.Order < jobs[j].Job.Order
	})

	return jobs
}
