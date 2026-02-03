package di

import (
	"context"

	fhmoscowrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhmoscow"
)

// FHMoscowTournamentRepository возвращает FHMoscow репозиторий турниров
func (c *Container) FHMoscowTournamentRepository(ctx context.Context) (*fhmoscowrepo.TournamentRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhmoscowrepo.NewTournamentRepository(db), nil
}

// FHMoscowTeamRepository возвращает FHMoscow репозиторий команд
func (c *Container) FHMoscowTeamRepository(ctx context.Context) (*fhmoscowrepo.TeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhmoscowrepo.NewTeamRepository(db), nil
}

// FHMoscowPlayerRepository возвращает FHMoscow репозиторий игроков
func (c *Container) FHMoscowPlayerRepository(ctx context.Context) (*fhmoscowrepo.PlayerRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhmoscowrepo.NewPlayerRepository(db), nil
}

// FHMoscowPlayerTeamRepository возвращает FHMoscow репозиторий связей
func (c *Container) FHMoscowPlayerTeamRepository(ctx context.Context) (*fhmoscowrepo.PlayerTeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhmoscowrepo.NewPlayerTeamRepository(db), nil
}

// FHMoscowPlayerStatisticsRepository возвращает FHMoscow репозиторий статистики
func (c *Container) FHMoscowPlayerStatisticsRepository(ctx context.Context) (*fhmoscowrepo.PlayerStatisticsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhmoscowrepo.NewPlayerStatisticsRepository(db), nil
}

// FHMoscowGoalieStatisticsRepository возвращает FHMoscow репозиторий статистики вратарей
func (c *Container) FHMoscowGoalieStatisticsRepository(ctx context.Context) (*fhmoscowrepo.GoalieStatisticsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhmoscowrepo.NewGoalieStatisticsRepository(db), nil
}
