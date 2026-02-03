package fhmoscow

import (
	"context"
	"fmt"

	fhmoscowrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhmoscow"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processPlayerSafe(ctx context.Context, tournamentID, teamID string, member dto.TeamMemberDTO) bool {
	err := o.processPlayer(ctx, tournamentID, teamID, member)
	if err != nil {
		logger.Warn(ctx, "Player processing failed",
			zap.String("id", member.PlayerID),
			zap.String("name", member.Name),
			zap.Error(err),
		)
		return false
	}
	return true
}

func (o *Orchestrator) processPlayer(ctx context.Context, tournamentID, teamID string, member dto.TeamMemberDTO) error {
	// Получаем профиль игрока
	profile, err := o.fetchPlayerProfile(ctx, member.PlayerID)
	if err != nil {
		logger.Debug(ctx, "Failed to fetch player profile, using basic data",
			zap.String("id", member.PlayerID),
			zap.Error(err),
		)
		// Используем базовые данные из состава команды
		profile = &dto.PlayerProfileDTO{
			ID:       member.PlayerID,
			FullName: member.Name,
			Position: member.Position,
		}
	}

	// Сохраняем игрока
	playerID, err := o.savePlayer(ctx, profile)
	if err != nil {
		return fmt.Errorf("save player: %w", err)
	}

	// Сохраняем связь игрок-команда
	if err := o.savePlayerTeam(ctx, playerID, teamID, tournamentID, member); err != nil {
		return fmt.Errorf("save player_team: %w", err)
	}

	// Сохраняем статистику (если есть в профиле)
	if err := o.savePlayerStatistics(ctx, playerID, teamID, tournamentID, profile); err != nil {
		logger.Debug(ctx, "Failed to save player statistics",
			zap.String("player_id", playerID),
			zap.Error(err),
		)
	}

	logger.Debug(ctx, "Player saved",
		zap.String("name", profile.FullName),
		zap.String("id", playerID),
	)

	return nil
}

func (o *Orchestrator) fetchPlayerProfile(ctx context.Context, playerID string) (*dto.PlayerProfileDTO, error) {
	path := fmt.Sprintf("/player/%s", playerID)
	html, err := o.client.GetHTML(path)
	if err != nil {
		return nil, fmt.Errorf("get player page: %w", err)
	}

	profile, err := parsing.ParsePlayerProfile(html, playerID)
	if err != nil {
		return nil, fmt.Errorf("parse profile: %w", err)
	}

	return profile, nil
}

func (o *Orchestrator) savePlayer(ctx context.Context, profile *dto.PlayerProfileDTO) (string, error) {
	profileURL := fmt.Sprintf("/player/%s", profile.ID)

	player := &fhmoscowrepo.Player{
		ExternalID: profile.ID,
		FullName:   profile.FullName,
		ProfileURL: &profileURL,
		BirthDate:  profile.BirthDate,
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

	return o.playerRepo.Upsert(ctx, player)
}

func (o *Orchestrator) savePlayerTeam(ctx context.Context, playerID, teamID, tournamentID string, member dto.TeamMemberDTO) error {
	pt := &fhmoscowrepo.PlayerTeam{
		PlayerID:     playerID,
		TeamID:       teamID,
		TournamentID: tournamentID,
	}

	if member.Number > 0 {
		pt.Number = &member.Number
	}
	if member.Position != "" {
		pt.Position = &member.Position
	}

	return o.playerTeamRepo.Upsert(ctx, pt)
}

func (o *Orchestrator) savePlayerStatistics(ctx context.Context, playerID, teamID, tournamentID string, profile *dto.PlayerProfileDTO) error {
	// Ищем статистику для текущего турнира/команды
	for _, stats := range profile.Stats {
		// Сохраняем статистику полевого игрока
		ps := &fhmoscowrepo.PlayerStatistics{
			PlayerID:       playerID,
			TeamID:         teamID,
			TournamentID:   tournamentID,
			Games:          stats.Games,
			Goals:          stats.Goals,
			Assists:        stats.Assists,
			Points:         stats.Points,
			PenaltyMinutes: stats.PenaltyMinutes,
		}

		if err := o.playerStatisticsRepo.Upsert(ctx, ps); err != nil {
			return err
		}

		// Берем только первую запись статистики (обычно текущий сезон)
		break
	}

	return nil
}

func (o *Orchestrator) saveGoalieStatistics(ctx context.Context, playerID, teamID, tournamentID string, stats dto.GoalieStatsDTO) error {
	gs := &fhmoscowrepo.GoalieStatistics{
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
