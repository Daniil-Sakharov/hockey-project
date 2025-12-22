package di

import (
	"context"

	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
)

// FHSPBTournamentRepository возвращает FHSPB репозиторий турниров
func (c *Container) FHSPBTournamentRepository(ctx context.Context) (*fhspbrepo.TournamentRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhspbrepo.NewTournamentRepository(db), nil
}

// FHSPBTeamRepository возвращает FHSPB репозиторий команд
func (c *Container) FHSPBTeamRepository(ctx context.Context) (*fhspbrepo.TeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhspbrepo.NewTeamRepository(db), nil
}

// FHSPBPlayerRepository возвращает FHSPB репозиторий игроков
func (c *Container) FHSPBPlayerRepository(ctx context.Context) (*fhspbrepo.PlayerRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhspbrepo.NewPlayerRepository(db), nil
}

// FHSPBPlayerTeamRepository возвращает FHSPB репозиторий связей
func (c *Container) FHSPBPlayerTeamRepository(ctx context.Context) (*fhspbrepo.PlayerTeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhspbrepo.NewPlayerTeamRepository(db), nil
}

// FHSPBPlayerStatisticsRepository возвращает FHSPB репозиторий статистики
func (c *Container) FHSPBPlayerStatisticsRepository(ctx context.Context) (*fhspbrepo.PlayerStatisticsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhspbrepo.NewPlayerStatisticsRepository(db), nil
}

// FHSPBGoalieStatisticsRepository возвращает FHSPB репозиторий статистики вратарей
func (c *Container) FHSPBGoalieStatisticsRepository(ctx context.Context) (*fhspbrepo.GoalieStatisticsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return fhspbrepo.NewGoalieStatisticsRepository(db), nil
}
