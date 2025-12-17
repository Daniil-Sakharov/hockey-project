package profile

import (
	"context"
	"fmt"

	domainPlayer "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	domainPlayerStats "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	domainPlayerTeam "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	domainTeam "github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

// Service сервис для работы с профилями игроков
type Service struct {
	playerRepo     domainPlayer.Repository
	statsRepo      domainPlayerStats.Repository
	playerTeamRepo domainPlayerTeam.Repository
	teamRepo       domainTeam.Repository
}

// NewService создает новый сервис профилей
func NewService(
	playerRepo domainPlayer.Repository,
	statsRepo domainPlayerStats.Repository,
	playerTeamRepo domainPlayerTeam.Repository,
	teamRepo domainTeam.Repository,
) *Service {
	return &Service{
		playerRepo:     playerRepo,
		statsRepo:      statsRepo,
		playerTeamRepo: playerTeamRepo,
		teamRepo:       teamRepo,
	}
}

// GetPlayerProfile получает полный профиль игрока
func (s *Service) GetPlayerProfile(ctx context.Context, playerID string) (*domainPlayer.Profile, error) {
	// 1. Получаем базовую информацию игрока
	player, err := s.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	// 2. Получаем команду и регион игрока
	team, region, err := s.getPlayerTeamAndRegion(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team and region: %w", err)
	}

	// 3. Формируем базовую информацию
	basicInfo := domainPlayer.PlayerBasicInfo{
		ID:        player.ID,
		Name:      player.Name,
		BirthYear: player.BirthDate.Year(),
		Position:  player.Position,
		Height:    player.Height,
		Weight:    player.Weight,
		Team:      team,
		Region:    region,
	}

	// 4. Получаем статистику за всё время
	allTimeStats, err := s.getAllTimeStats(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all time stats: %w", err)
	}

	// Если нет статистики вообще - возвращаем только базовую информацию
	if allTimeStats == nil {
		return &domainPlayer.Profile{
			BasicInfo:          basicInfo,
			AllTimeStats:       nil,
			CurrentSeasonStats: nil,
			RecentTournaments:  []domainPlayer.TournamentStats{},
		}, nil
	}

	// 5. Определяем текущий сезон и получаем статистику за него
	currentSeason := getCurrentSeason()
	seasonStats, err := s.getSeasonStats(ctx, playerID, currentSeason)
	if err != nil {
		return nil, fmt.Errorf("failed to get season stats: %w", err)
	}

	// Если нет статистики за текущий сезон - пытаемся взять предыдущий
	if seasonStats == nil {
		previousSeason := getPreviousSeason(currentSeason)
		seasonStats, err = s.getSeasonStats(ctx, playerID, previousSeason)
		if err != nil {
			return nil, fmt.Errorf("failed to get previous season stats: %w", err)
		}
	}

	// 6. Получаем последние 5 турниров
	recentTournaments, err := s.getRecentTournaments(ctx, playerID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tournaments: %w", err)
	}

	return &domainPlayer.Profile{
		BasicInfo:          basicInfo,
		AllTimeStats:       allTimeStats,
		CurrentSeasonStats: seasonStats,
		RecentTournaments:  recentTournaments,
	}, nil
}
