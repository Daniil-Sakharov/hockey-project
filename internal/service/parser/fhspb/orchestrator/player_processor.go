package orchestrator

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// processPlayerSafe обрабатывает игрока с логированием ошибок
func (o *Orchestrator) processPlayerSafe(ctx context.Context, pURL dto.PlayerURLDTO) bool {
	err := o.processPlayer(ctx, pURL)
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
func (o *Orchestrator) processPlayer(ctx context.Context, pURL dto.PlayerURLDTO) error {
	exists, err := o.playerRepo.ExistsByExternalID(ctx, pURL.PlayerID, player.SourceFHSPB)
	if err != nil {
		return fmt.Errorf("check exists: %w", err)
	}

	if exists {
		return nil
	}

	playerDTO, err := o.client.GetPlayer(pURL.TournamentID, pURL.TeamID, pURL.PlayerID)
	if err != nil {
		return fmt.Errorf("get player: %w", err)
	}

	p, err := convertToPlayer(playerDTO)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	if err := o.playerRepo.Upsert(ctx, p); err != nil {
		return fmt.Errorf("upsert: %w", err)
	}

	logger.Debug(ctx, "💾 Player saved", zap.String("name", p.Name))

	return nil
}
