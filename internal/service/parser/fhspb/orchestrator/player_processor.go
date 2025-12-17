package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	fhspbRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// processPlayerSafe обрабатывает игрока с логированием ошибок
func (o *Orchestrator) processPlayerSafe(ctx context.Context, teamID, tournamentID string, pURL dto.PlayerURLDTO) bool {
	err := o.processPlayer(ctx, teamID, tournamentID, pURL)
	if err != nil {
		logger.Warn(ctx, "⚠️ Player failed",
			zap.String("id", pURL.PlayerID),
			zap.Error(err),
		)
		return false
	}
	return true
}

// processPlayer обрабатывает и сохраняет одного игрока
func (o *Orchestrator) processPlayer(ctx context.Context, teamID, tournamentID string, pURL dto.PlayerURLDTO) error {
	playerDTO, err := o.client.GetPlayer(pURL.TournamentID, pURL.TeamID, pURL.PlayerID)
	if err != nil {
		return fmt.Errorf("get player: %w", err)
	}

	// Сохраняем/обновляем профиль игрока
	player := &fhspbRepo.Player{
		ExternalID: pURL.PlayerID,
		FullName:   playerDTO.FullName,
	}
	if playerDTO.BirthDate != "" {
		if t, err := time.Parse("02.01.2006", playerDTO.BirthDate); err == nil {
			player.BirthDate = &t
		}
	}
	if playerDTO.BirthPlace != "" {
		player.BirthPlace = &playerDTO.BirthPlace
	}

	playerID, err := o.playerRepo.Upsert(ctx, player)
	if err != nil {
		return fmt.Errorf("upsert player: %w", err)
	}

	// Сохраняем связь игрок-команда
	playerTeam := &fhspbRepo.PlayerTeam{
		PlayerID:     playerID,
		TeamID:       teamID,
		TournamentID: tournamentID,
	}
	if playerDTO.Number > 0 {
		playerTeam.Number = &playerDTO.Number
	}
	if playerDTO.Role != "" {
		playerTeam.Role = &playerDTO.Role
	}
	if playerDTO.Position != "" {
		playerTeam.Position = &playerDTO.Position
	}

	if err := o.playerTeamRepo.Upsert(ctx, playerTeam); err != nil {
		return fmt.Errorf("upsert player_team: %w", err)
	}

	logger.Info(ctx, "👤 Player saved", zap.String("name", playerDTO.FullName), zap.String("id", playerID))
	return nil
}
