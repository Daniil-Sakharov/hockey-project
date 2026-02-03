package mihf

import (
	"context"
	"fmt"
	"time"

	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/retry"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processPlayerStatsSafe(ctx context.Context, tournamentID, teamID string, birthYear int, stats dto.PlayerStatsDTO) bool {
	err := o.processPlayerStats(ctx, tournamentID, teamID, birthYear, stats)
	if err != nil {
		logger.Warn(ctx, "Player stats failed",
			zap.String("id", stats.ID),
			zap.String("name", stats.Name),
			zap.Error(err),
		)

		if o.config.RetryEnabled() {
			retryErr := o.retryManager.AddFailedJob(ctx, retry.JobTypePlayer, "mihf", stats.ID, stats.ProfileURL, err)
			if retryErr != nil {
				logger.Error(ctx, "Failed to add retry job", zap.Error(retryErr))
			}
		}

		return false
	}
	return true
}

func (o *Orchestrator) processPlayerStats(ctx context.Context, tournamentID, teamID string, birthYear int, stats dto.PlayerStatsDTO) error {
	// Получаем профиль игрока для антропометрии
	profile, err := o.fetchPlayerProfile(ctx, stats.ID)
	if err != nil {
		logger.Debug(ctx, "Failed to fetch player profile, saving basic data",
			zap.String("id", stats.ID),
			zap.Error(err),
		)
		// Продолжаем с базовыми данными
		profile = &dto.PlayerProfileDTO{
			ID:       stats.ID,
			FullName: stats.Name,
		}
	}

	// Сохраняем игрока
	playerID, err := o.savePlayer(ctx, profile, stats, birthYear)
	if err != nil {
		return fmt.Errorf("save player: %w", err)
	}

	// Сохраняем связь игрок-команда
	if err := o.savePlayerTeam(ctx, playerID, teamID, tournamentID, stats); err != nil {
		return fmt.Errorf("save player_team: %w", err)
	}

	// Сохраняем статистику
	if err := o.savePlayerStatistics(ctx, playerID, teamID, tournamentID, stats); err != nil {
		return fmt.Errorf("save player_statistics: %w", err)
	}

	logger.Debug(ctx, "Player saved", zap.String("name", stats.Name), zap.String("id", playerID))
	return nil
}

func (o *Orchestrator) processGoalieStatsSafe(ctx context.Context, tournamentID, teamID string, birthYear int, stats dto.GoalieStatsDTO) bool {
	err := o.processGoalieStats(ctx, tournamentID, teamID, birthYear, stats)
	if err != nil {
		logger.Warn(ctx, "Goalie stats failed",
			zap.String("id", stats.ID),
			zap.String("name", stats.Name),
			zap.Error(err),
		)
		return false
	}
	return true
}

func (o *Orchestrator) processGoalieStats(ctx context.Context, tournamentID, teamID string, birthYear int, stats dto.GoalieStatsDTO) error {
	// Получаем профиль
	profile, err := o.fetchPlayerProfile(ctx, stats.ID)
	if err != nil {
		profile = &dto.PlayerProfileDTO{
			ID:       stats.ID,
			FullName: stats.Name,
			Position: "В",
		}
	}

	// Устанавливаем позицию вратаря
	if profile.Position == "" {
		profile.Position = "В"
	}

	// Сохраняем игрока
	playerID, err := o.saveGoalieAsPlayer(ctx, profile, stats, birthYear)
	if err != nil {
		return fmt.Errorf("save goalie: %w", err)
	}

	// Сохраняем связь игрок-команда
	if err := o.saveGoalieTeam(ctx, playerID, teamID, tournamentID, stats); err != nil {
		return fmt.Errorf("save player_team: %w", err)
	}

	// Сохраняем статистику вратаря
	if err := o.saveGoalieStatistics(ctx, playerID, teamID, tournamentID, stats); err != nil {
		return fmt.Errorf("save goalie_statistics: %w", err)
	}

	logger.Debug(ctx, "Goalie saved", zap.String("name", stats.Name), zap.String("id", playerID))
	return nil
}

func (o *Orchestrator) fetchPlayerProfile(ctx context.Context, playerID string) (*dto.PlayerProfileDTO, error) {
	url := fmt.Sprintf("/players/info/%s", playerID)
	html, err := o.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get profile page: %w", err)
	}

	profile, err := parsing.ParsePlayerProfile(html, playerID)
	if err != nil {
		return nil, fmt.Errorf("parse profile: %w", err)
	}

	return profile, nil
}

func (o *Orchestrator) savePlayer(ctx context.Context, profile *dto.PlayerProfileDTO, stats dto.PlayerStatsDTO, birthYear int) (string, error) {
	player := &mihfrepo.Player{
		ExternalID: stats.ID,
		FullName:   profile.FullName,
	}

	if player.FullName == "" {
		player.FullName = stats.Name
	}

	if profile.BirthDate != nil {
		player.BirthDate = profile.BirthDate
	} else if birthYear > 0 {
		// Fallback: используем год рождения турнира
		fallbackDate := time.Date(birthYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		player.BirthDate = &fallbackDate
	}

	// Позиция: приоритет профилю, затем таблице статистики
	if profile.Position != "" {
		player.Position = &profile.Position
	} else if stats.Position != "" {
		player.Position = &stats.Position
	}

	if profile.Height > 0 {
		player.Height = &profile.Height
	}
	if profile.Weight > 0 {
		player.Weight = &profile.Weight
	}
	if profile.Handedness != "" {
		player.Handedness = &profile.Handedness
	}
	if profile.Citizenship != "" {
		player.Citizenship = &profile.Citizenship
	}
	if stats.ProfileURL != "" {
		player.ProfileURL = &stats.ProfileURL
	}

	return o.playerRepo.Upsert(ctx, player)
}

func (o *Orchestrator) saveGoalieAsPlayer(ctx context.Context, profile *dto.PlayerProfileDTO, stats dto.GoalieStatsDTO, birthYear int) (string, error) {
	player := &mihfrepo.Player{
		ExternalID: stats.ID,
		FullName:   profile.FullName,
	}

	if player.FullName == "" {
		player.FullName = stats.Name
	}

	if profile.BirthDate != nil {
		player.BirthDate = profile.BirthDate
	} else if birthYear > 0 {
		// Fallback: используем год рождения турнира
		fallbackDate := time.Date(birthYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		player.BirthDate = &fallbackDate
	}

	if profile.Position != "" {
		player.Position = &profile.Position
	}
	if profile.Height > 0 {
		player.Height = &profile.Height
	}
	if profile.Weight > 0 {
		player.Weight = &profile.Weight
	}
	if profile.Handedness != "" {
		player.Handedness = &profile.Handedness
	}
	if profile.Citizenship != "" {
		player.Citizenship = &profile.Citizenship
	}
	if stats.ProfileURL != "" {
		player.ProfileURL = &stats.ProfileURL
	}

	return o.playerRepo.Upsert(ctx, player)
}

func (o *Orchestrator) savePlayerTeam(ctx context.Context, playerID, teamID, tournamentID string, stats dto.PlayerStatsDTO) error {
	pt := &mihfrepo.PlayerTeam{
		PlayerID:     playerID,
		TeamID:       teamID,
		TournamentID: tournamentID,
	}

	if stats.Number != "" {
		num := parseNumber(stats.Number)
		if num > 0 {
			pt.Number = &num
		}
	}

	// Передаем позицию из таблицы (З/Н)
	if stats.Position != "" {
		pt.Position = &stats.Position
	}

	return o.playerTeamRepo.Upsert(ctx, pt)
}

func (o *Orchestrator) saveGoalieTeam(ctx context.Context, playerID, teamID, tournamentID string, stats dto.GoalieStatsDTO) error {
	pt := &mihfrepo.PlayerTeam{
		PlayerID:     playerID,
		TeamID:       teamID,
		TournamentID: tournamentID,
	}

	if stats.Number != "" {
		num := parseNumber(stats.Number)
		if num > 0 {
			pt.Number = &num
		}
	}

	position := "В"
	pt.Position = &position

	return o.playerTeamRepo.Upsert(ctx, pt)
}

func (o *Orchestrator) savePlayerStatistics(ctx context.Context, playerID, teamID, tournamentID string, stats dto.PlayerStatsDTO) error {
	ps := &mihfrepo.PlayerStatistics{
		PlayerID:          playerID,
		TeamID:            teamID,
		TournamentID:      tournamentID,
		Games:             stats.Games,
		Goals:             stats.Goals,
		Assists:           stats.Assists,
		Points:            stats.Points,
		PenaltyMinutes:    stats.PenaltyMinutes,
		GoalsPowerPlay:    &stats.GoalsPowerPlay,
		GoalsShortHanded:  &stats.GoalsShortHanded,
		GoalsEvenStrength: &stats.GoalsEvenStrength,
	}

	return o.playerStatisticsRepo.Upsert(ctx, ps)
}

func (o *Orchestrator) saveGoalieStatistics(ctx context.Context, playerID, teamID, tournamentID string, stats dto.GoalieStatsDTO) error {
	gs := &mihfrepo.GoalieStatistics{
		PlayerID:       playerID,
		TeamID:         teamID,
		TournamentID:   tournamentID,
		Games:          stats.Games,
		GoalsAgainst:   stats.GoalsAgainst,
		SavePercentage: &stats.SavePercentage,
		MinutesPlayed:  &stats.MinutesPlayed,
	}

	return o.goalieStatisticsRepo.Upsert(ctx, gs)
}

func parseNumber(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
