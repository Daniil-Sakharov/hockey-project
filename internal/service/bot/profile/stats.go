package profile

import (
	"context"

	domainPlayer "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// getAllTimeStats преобразует агрегированную статистику в доменную модель
func (s *Service) getAllTimeStats(ctx context.Context, playerID string) (*domainPlayer.PlayerStats, error) {
	stats, err := s.statsRepo.GetAllTimeStats(ctx, playerID)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, nil
	}

	return &domainPlayer.PlayerStats{
		Tournaments:      stats.TournamentsCount,
		Games:            stats.Games,
		Goals:            stats.Goals,
		Assists:          stats.Assists,
		Points:           stats.Points,
		PlusMinus:        stats.PlusMinus,
		Penalties:        stats.PenaltyMinutes,
		GoalsPerGame:     calculateAverage(stats.Goals, stats.Games),
		AssistsPerGame:   calculateAverage(stats.Assists, stats.Games),
		PointsPerGame:    calculateAverage(stats.Points, stats.Games),
		PenaltiesPerGame: calculateAverage(stats.PenaltyMinutes, stats.Games),
		HatTricks:        stats.HatTricks,
		GameWinningGoals: stats.GameWinningGoals,
	}, nil
}

// getSeasonStats преобразует статистику сезона в доменную модель
func (s *Service) getSeasonStats(ctx context.Context, playerID, season string) (*domainPlayer.SeasonStats, error) {
	stats, err := s.statsRepo.GetSeasonStats(ctx, playerID, season)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, nil
	}

	return &domainPlayer.SeasonStats{
		Season:           season,
		Tournaments:      stats.TournamentsCount,
		Games:            stats.Games,
		Goals:            stats.Goals,
		Assists:          stats.Assists,
		Points:           stats.Points,
		PlusMinus:        stats.PlusMinus,
		Penalties:        stats.PenaltyMinutes,
		GoalsPerGame:     calculateAverage(stats.Goals, stats.Games),
		AssistsPerGame:   calculateAverage(stats.Assists, stats.Games),
		PointsPerGame:    calculateAverage(stats.Points, stats.Games),
		PenaltiesPerGame: calculateAverage(stats.PenaltyMinutes, stats.Games),
		HatTricks:        stats.HatTricks,
		GameWinningGoals: stats.GameWinningGoals,
	}, nil
}
