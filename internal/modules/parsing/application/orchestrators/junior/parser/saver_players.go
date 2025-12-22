package parser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// SavePlayers парсит игроков одной команды и сохраняет в БД
func (s *orchestratorService) SavePlayers(ctx context.Context, teamURL, teamID, tournamentID string, t *entities.Tournament) error {
	logger.Info(ctx, fmt.Sprintf("  🏒 Parsing team: %s", teamURL))

	playersDTO, err := s.juniorService.ParsePlayers(ctx, teamURL)
	if err != nil {
		return fmt.Errorf("failed to parse players: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("  ✅ Found %d players in HTML", len(playersDTO)))

	if len(playersDTO) == 0 {
		logger.Warn(ctx, "  ⚠️  NO PLAYERS FOUND!")
		return nil
	}

	currentSeason := ""
	if t != nil {
		currentSeason = t.Season
	}

	savedCount, updatedCount, existingCount, skippedCount, skippedTooOld := 0, 0, 0, 0, 0
	playerIDs := make([]string, 0, len(playersDTO))

	for i, dto := range playersDTO {
		logger.Info(ctx, fmt.Sprintf("    [%d/%d] Processing: %s (%s)", i+1, len(playersDTO), dto.Name, dto.ProfileURL))

		p, err := s.convertPlayerDTO(dto, currentSeason)
		if err != nil {
			if strings.Contains(err.Error(), "too old") {
				skippedTooOld++
			} else {
				skippedCount++
			}
			continue
		}

		existing, err := s.playerRepo.GetByProfileURL(ctx, p.ProfileURL)
		if err != nil {
			return fmt.Errorf("failed to check existing player: %w", err)
		}

		if existing == nil {
			logger.Info(ctx, fmt.Sprintf("    ✅ Creating NEW player: %s (ID: %s)", p.Name, p.ID))
			if err := s.playerRepo.Create(ctx, p); err != nil {
				return fmt.Errorf("failed to create player %s: %w", p.Name, err)
			}
			savedCount++
			playerIDs = append(playerIDs, p.ID)
		} else {
			updated := s.updatePlayerIfNeeded(ctx, existing, p, currentSeason)
			if updated {
				updatedCount++
			}
			existingCount++
			playerIDs = append(playerIDs, existing.ID)
		}
	}

	if len(playerIDs) > 0 {
		if err := s.CreatePlayerTeamLinksBatch(ctx, playerIDs, teamID, tournamentID, t); err != nil {
			return fmt.Errorf("failed to create player_team links: %w", err)
		}
	}

	logger.Info(ctx, fmt.Sprintf("  📊 ИТОГО: новых=%d, обновлено=%d, существующих=%d, пропущено=%d, слишком старых=%d",
		savedCount, updatedCount, existingCount, skippedCount, skippedTooOld))
	return nil
}

// updatePlayerIfNeeded обновляет данные игрока если текущий сезон новее
func (s *orchestratorService) updatePlayerIfNeeded(ctx context.Context, existing, p *entities.Player, currentSeason string) bool {
	existingDataSeason := ""
	if existing.DataSeason != nil {
		existingDataSeason = *existing.DataSeason
	}

	if !entities.IsNewerSeason(currentSeason, existingDataSeason) {
		return false
	}

	needUpdate := false

	if p.Height != nil && (existing.Height == nil || *existing.Height != *p.Height) {
		existing.Height = p.Height
		needUpdate = true
	}
	if p.Weight != nil && (existing.Weight == nil || *existing.Weight != *p.Weight) {
		existing.Weight = p.Weight
		needUpdate = true
	}
	if p.Handedness != nil && (existing.Handedness == nil || *existing.Handedness != *p.Handedness) {
		existing.Handedness = p.Handedness
		needUpdate = true
	}
	if p.Position != "" && existing.Position != p.Position {
		existing.Position = p.Position
		needUpdate = true
	}

	if needUpdate {
		existing.DataSeason = &currentSeason
		existing.UpdatedAt = time.Now()

		if err := s.playerRepo.Update(ctx, existing); err != nil {
			logger.Warn(ctx, fmt.Sprintf("    ⚠️  Failed to update player: %v", err))
			return false
		}
		return true
	}

	return false
}
