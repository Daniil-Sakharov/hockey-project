package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type StatisticsPostgres struct {
	db *sqlx.DB
}

func NewStatisticsPostgres(db *sqlx.DB) *StatisticsPostgres {
	return &StatisticsPostgres{db: db}
}

// CreateBatch создает множество статистик за один запрос с FK validation
func (r *StatisticsPostgres) CreateBatch(ctx context.Context, stats []*entities.PlayerStatistic) (int, error) {
	if len(stats) == 0 {
		return 0, nil
	}

	validStats, err := r.filterValidStats(ctx, stats)
	if err != nil {
		return 0, err
	}

	if len(validStats) == 0 {
		return 0, nil
	}

	return r.insertBatch(ctx, validStats)
}

// filterValidStats фильтрует статистики по существующим FK
func (r *StatisticsPostgres) filterValidStats(ctx context.Context, stats []*entities.PlayerStatistic) ([]*entities.PlayerStatistic, error) {
	playerIDsMap := make(map[string]bool)
	teamIDsMap := make(map[string]bool)

	for _, stat := range stats {
		playerIDsMap[stat.PlayerID] = true
		teamIDsMap[stat.TeamID] = true
	}

	playerIDsList := make([]string, 0, len(playerIDsMap))
	for id := range playerIDsMap {
		playerIDsList = append(playerIDsList, id)
	}

	teamIDsList := make([]string, 0, len(teamIDsMap))
	for id := range teamIDsMap {
		teamIDsList = append(teamIDsList, id)
	}

	var existingPlayers []string
	if len(playerIDsList) > 0 {
		err := r.db.SelectContext(ctx, &existingPlayers, `SELECT id FROM players WHERE id = ANY($1)`, pq.Array(playerIDsList))
		if err != nil {
			return nil, fmt.Errorf("failed to check existing players: %w", err)
		}
	}

	var existingTeams []string
	if len(teamIDsList) > 0 {
		err := r.db.SelectContext(ctx, &existingTeams, `SELECT id FROM teams WHERE id = ANY($1)`, pq.Array(teamIDsList))
		if err != nil {
			return nil, fmt.Errorf("failed to check existing teams: %w", err)
		}
	}

	existingPlayersMap := make(map[string]bool)
	for _, id := range existingPlayers {
		existingPlayersMap[id] = true
	}

	existingTeamsMap := make(map[string]bool)
	for _, id := range existingTeams {
		existingTeamsMap[id] = true
	}

	validStats := make([]*entities.PlayerStatistic, 0, len(stats))
	for _, stat := range stats {
		if existingPlayersMap[stat.PlayerID] && existingTeamsMap[stat.TeamID] {
			validStats = append(validStats, stat)
		}
	}

	return validStats, nil
}
