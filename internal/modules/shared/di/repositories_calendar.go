package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/repositories"
	parsingRepos "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories"
)

// MatchRepository возвращает репозиторий матчей
func (c *Container) MatchRepository(ctx context.Context) (repositories.MatchRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewMatchPostgres(db), nil
}

// MatchEventRepository возвращает репозиторий событий матчей
func (c *Container) MatchEventRepository(ctx context.Context) (repositories.MatchEventRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewMatchEventPostgres(db), nil
}

// MatchLineupRepository возвращает репозиторий составов матчей
func (c *Container) MatchLineupRepository(ctx context.Context) (repositories.MatchLineupRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewMatchLineupPostgres(db), nil
}

// StandingRepository возвращает репозиторий турнирных таблиц
func (c *Container) StandingRepository(ctx context.Context) (repositories.StandingRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewStandingPostgres(db), nil
}

// MatchTeamStatsRepository возвращает репозиторий статистики команд за матч
func (c *Container) MatchTeamStatsRepository(ctx context.Context) (repositories.MatchTeamStatsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewMatchTeamStatsPostgres(db), nil
}
