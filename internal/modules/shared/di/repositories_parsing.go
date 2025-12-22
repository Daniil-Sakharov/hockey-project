package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/repositories"
	parsingRepos "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories"
)

// ParsingPlayerRepository возвращает репозиторий игроков из parsing модуля
func (c *Container) ParsingPlayerRepository(ctx context.Context) (repositories.PlayerRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewPlayerPostgres(db), nil
}

// ParsingTeamRepository возвращает репозиторий команд из parsing модуля
func (c *Container) ParsingTeamRepository(ctx context.Context) (repositories.TeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewTeamPostgres(db), nil
}

// ParsingTournamentRepository возвращает репозиторий турниров из parsing модуля
func (c *Container) ParsingTournamentRepository(ctx context.Context) (repositories.TournamentRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewTournamentPostgres(db), nil
}

// ParsingPlayerTeamRepository возвращает репозиторий связей игрок-команда из parsing модуля
func (c *Container) ParsingPlayerTeamRepository(ctx context.Context) (repositories.PlayerTeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewPlayerTeamRepository(db), nil
}

// ParsingPlayerStatisticsRepository возвращает репозиторий статистики из parsing модуля
func (c *Container) ParsingPlayerStatisticsRepository(ctx context.Context) (repositories.PlayerStatisticsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewStatisticsPostgres(db), nil
}

// ParsingGoalieStatisticsRepository возвращает репозиторий статистики вратарей из parsing модуля
func (c *Container) ParsingGoalieStatisticsRepository(ctx context.Context) (*parsingRepos.GoalieStatisticsPostgres, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return parsingRepos.NewGoalieStatisticsPostgres(db), nil
}
