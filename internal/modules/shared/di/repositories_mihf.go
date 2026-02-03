package di

import (
	"context"

	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
)

// MIHFTournamentRepository возвращает MIHF репозиторий турниров
func (c *Container) MIHFTournamentRepository(ctx context.Context) (*mihfrepo.TournamentRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return mihfrepo.NewTournamentRepository(db), nil
}

// MIHFTeamRepository возвращает MIHF репозиторий команд
func (c *Container) MIHFTeamRepository(ctx context.Context) (*mihfrepo.TeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return mihfrepo.NewTeamRepository(db), nil
}

// MIHFPlayerRepository возвращает MIHF репозиторий игроков
func (c *Container) MIHFPlayerRepository(ctx context.Context) (*mihfrepo.PlayerRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return mihfrepo.NewPlayerRepository(db), nil
}

// MIHFPlayerTeamRepository возвращает MIHF репозиторий связей
func (c *Container) MIHFPlayerTeamRepository(ctx context.Context) (*mihfrepo.PlayerTeamRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return mihfrepo.NewPlayerTeamRepository(db), nil
}

// MIHFPlayerStatisticsRepository возвращает MIHF репозиторий статистики
func (c *Container) MIHFPlayerStatisticsRepository(ctx context.Context) (*mihfrepo.PlayerStatisticsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return mihfrepo.NewPlayerStatisticsRepository(db), nil
}

// MIHFGoalieStatisticsRepository возвращает MIHF репозиторий статистики вратарей
func (c *Container) MIHFGoalieStatisticsRepository(ctx context.Context) (*mihfrepo.GoalieStatisticsRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return mihfrepo.NewGoalieStatisticsRepository(db), nil
}
