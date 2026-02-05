package parser

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// SavePlayers парсит игроков одной команды и сохраняет в БД
func (s *orchestratorService) SavePlayers(ctx context.Context, domain, teamURL, teamID, tournamentID string, t *entities.Tournament, birthYear *int, groupName *string) error {
	logger.Info(ctx, fmt.Sprintf("  🏒 Parsing team: %s", teamURL))

	playersDTO, err := s.juniorService.ParsePlayers(ctx, domain, teamURL)
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
	playerCtx := make(map[string]playerLinkContext) // playerID → context

	for i, dto := range playersDTO {
		logger.Info(ctx, fmt.Sprintf("    [%d/%d] Processing: %s (%s)", i+1, len(playersDTO), dto.Name, dto.ProfileURL))

		// Дополняем данные из профиля если чего-то не хватает
		if s.needsProfileFetch(&dto) {
			s.enrichFromProfile(ctx, domain, &dto)
		}

		p, err := s.convertPlayerDTO(dto, currentSeason, domain)
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

		var playerID string
		if existing == nil {
			logger.Info(ctx, fmt.Sprintf("    ✅ Creating NEW player: %s (ID: %s)", p.Name, p.ID))
			if err := s.playerRepo.Create(ctx, p); err != nil {
				return fmt.Errorf("failed to create player %s: %w", p.Name, err)
			}
			savedCount++
			playerID = p.ID
		} else {
			updated := s.updatePlayerIfNeeded(ctx, existing, p, currentSeason)
			if updated {
				updatedCount++
			}
			existingCount++
			playerID = existing.ID
		}

		playerIDs = append(playerIDs, playerID)
		pctx := playerLinkContext{PhotoURL: dto.PhotoURL}
		if dto.Number != "" {
			if num, err := strconv.Atoi(dto.Number); err == nil {
				pctx.JerseyNumber = &num
			}
		}
		if dto.Height != "" {
			if h, err := strconv.Atoi(strings.TrimSpace(dto.Height)); err == nil {
				pctx.Height = &h
			}
		}
		if dto.Weight != "" {
			if w, err := strconv.Atoi(strings.TrimSpace(dto.Weight)); err == nil {
				pctx.Weight = &w
			}
		}
		playerCtx[playerID] = pctx
	}

	if len(playerIDs) > 0 {
		if err := s.CreatePlayerTeamLinksBatch(ctx, playerIDs, teamID, tournamentID, t, birthYear, groupName, playerCtx); err != nil {
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
	if p.Citizenship != nil && (existing.Citizenship == nil || *existing.Citizenship != *p.Citizenship) {
		existing.Citizenship = p.Citizenship
		needUpdate = true
	}
	if p.PhotoURL != nil && (existing.PhotoURL == nil || *existing.PhotoURL != *p.PhotoURL) {
		existing.PhotoURL = p.PhotoURL
		needUpdate = true
	}
	if p.Domain != nil && (existing.Domain == nil || *existing.Domain != *p.Domain) {
		existing.Domain = p.Domain
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

// needsProfileFetch проверяет нужно ли загружать профиль игрока
func (s *orchestratorService) needsProfileFetch(dto *types.PlayerDTO) bool {
	return dto.Position == "" || dto.Citizenship == "" || dto.PhotoURL == ""
}

// enrichFromProfile дополняет данные игрока из профиля
func (s *orchestratorService) enrichFromProfile(ctx context.Context, domain string, dto *types.PlayerDTO) {
	if dto.ProfileURL == "" {
		return
	}

	profile, err := s.juniorService.ParsePlayerProfile(ctx, domain, dto.ProfileURL)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("    ⚠️  Failed to fetch profile %s: %v", dto.ProfileURL, err))
		return
	}

	if profile == nil {
		return
	}

	// Дополняем только отсутствующие поля
	if dto.Position == "" && profile.Position != "" {
		dto.Position = profile.Position
		logger.Info(ctx, fmt.Sprintf("    📋 Got position from profile: %s", profile.Position))
	}
	if dto.Height == "" && profile.Height != "" {
		dto.Height = profile.Height
	}
	if dto.Weight == "" && profile.Weight != "" {
		dto.Weight = profile.Weight
	}
	if dto.Handedness == "" && profile.Handedness != "" {
		dto.Handedness = profile.Handedness
	}
	if dto.Citizenship == "" && profile.Citizenship != "" {
		dto.Citizenship = profile.Citizenship
		logger.Info(ctx, fmt.Sprintf("    🌍 Got citizenship from profile: %s", profile.Citizenship))
	}
	if dto.PhotoURL == "" && profile.PhotoURL != "" {
		dto.PhotoURL = profile.PhotoURL
	}
}
