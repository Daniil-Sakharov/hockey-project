package tasks

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/workers/pool"
)

// ParsingTask задача парсинга
type ParsingTask struct {
	*pool.BaseTask
	URL    string
	Domain string
	Type   string // "tournament", "team", "player"
}

// NewParsingTask создает задачу парсинга
func NewParsingTask(id, url, domain, taskType string, handler func(ctx context.Context) (interface{}, error)) *ParsingTask {
	baseTask := pool.NewBaseTask(id, 5, handler) // Средний приоритет

	return &ParsingTask{
		BaseTask: baseTask,
		URL:      url,
		Domain:   domain,
		Type:     taskType,
	}
}

// TeamTask задача обработки команды
type TeamTask struct {
	*pool.BaseTask
	TeamID       string
	TournamentID string
}

// NewTeamTask создает задачу обработки команды
func NewTeamTask(teamID, tournamentID string, handler func(ctx context.Context) (interface{}, error)) *TeamTask {
	id := fmt.Sprintf("team-%s-%s", teamID, tournamentID)
	baseTask := pool.NewBaseTask(id, 7, handler) // Высокий приоритет

	return &TeamTask{
		BaseTask:     baseTask,
		TeamID:       teamID,
		TournamentID: tournamentID,
	}
}

// PlayerTask задача обработки игрока
type PlayerTask struct {
	*pool.BaseTask
	PlayerID string
	TeamID   string
}

// NewPlayerTask создает задачу обработки игрока
func NewPlayerTask(playerID, teamID string, handler func(ctx context.Context) (interface{}, error)) *PlayerTask {
	id := fmt.Sprintf("player-%s-%s", playerID, teamID)
	baseTask := pool.NewBaseTask(id, 3, handler) // Низкий приоритет

	return &PlayerTask{
		BaseTask: baseTask,
		PlayerID: playerID,
		TeamID:   teamID,
	}
}
