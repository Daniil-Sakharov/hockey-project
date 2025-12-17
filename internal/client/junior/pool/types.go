package pool

import (
	jrtypes "github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/types"
)

// SeasonTask задача для воркера (парсинг одного сезона)
type SeasonTask struct {
	Domain  string
	Season  string
	AjaxURL string
}

// WorkerPoolResult результат работы воркера
type WorkerPoolResult struct {
	Tournaments []jrtypes.TournamentDTO
	Error       error
	Task        SeasonTask
}
