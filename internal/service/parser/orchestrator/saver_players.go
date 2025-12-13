package orchestrator

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// SavePlayers парсит игроков одной команды и сохраняет в БД + создает связи player_teams (BATCH)
// Обновляет данные игрока (рост, вес, хват) если текущий турнир из более нового сезона
func (s *orchestratorService) SavePlayers(ctx context.Context, teamURL, teamID, tournamentID string, t *tournament.Tournament) error {
	logger.Info(ctx, fmt.Sprintf("  🏒 Parsing team: %s", teamURL))

	// Парсим игроков через JuniorService
	playersDTO, err := s.juniorService.ParsePlayers(ctx, teamURL)
	if err != nil {
		return fmt.Errorf("failed to parse players: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("  ✅ Found %d players in HTML", len(playersDTO)))

	if len(playersDTO) == 0 {
		logger.Warn(ctx, "  ⚠️  NO PLAYERS FOUND! Check HTML selectors or page structure")
		return nil
	}

	// Получаем сезон текущего турнира для обновления данных
	currentSeason := ""
	if t != nil {
		currentSeason = t.Season
	}

	// Конвертируем DTO в domain entity
	savedCount := 0
	skippedCount := 0
	skippedTooOld := 0 // Пропущено по возрасту (год рождения < MinBirthYear)
	existingCount := 0
	updatedCount := 0
	playerIDs := make([]string, 0, len(playersDTO)) // Для batch создания связей

	for i, dto := range playersDTO {
		logger.Info(ctx, fmt.Sprintf("    [%d/%d] Processing: %s (%s)", i+1, len(playersDTO), dto.Name, dto.ProfileURL))

		// Конвертируем DTO → Entity
		p, err := convertPlayerDTO(dto, currentSeason)
		if err != nil {
			// Проверяем причину ошибки - пропуск по возрасту или другая ошибка
			if strings.Contains(err.Error(), "too old") {
				logger.Info(ctx, fmt.Sprintf("    ⏭️  Skipped (too old): %s - %s", dto.Name, dto.BirthDate))
				skippedTooOld++
			} else {
				logger.Warn(ctx, fmt.Sprintf("    ⚠️  Skipped (conversion error): %v", err))
				skippedCount++
			}
			continue
		}

		// Проверяем существование по URL (дедупликация)
		existing, err := s.playerRepo.GetByProfileURL(ctx, p.ProfileURL)
		if err != nil {
			return fmt.Errorf("failed to check existing player: %w", err)
		}

		// Если игрока нет - создаем
		if existing == nil {
			logger.Info(ctx, fmt.Sprintf("    ✅ Creating NEW player: %s (ID: %s)", p.Name, p.ID))
			if err := s.playerRepo.Create(ctx, p); err != nil {
				return fmt.Errorf("failed to create player %s: %w", p.Name, err)
			}
			savedCount++
			playerIDs = append(playerIDs, p.ID)
		} else {
			// Игрок существует - проверяем нужно ли обновить данные
			existingDataSeason := ""
			if existing.DataSeason != nil {
				existingDataSeason = *existing.DataSeason
			}

			// Обновляем данные если текущий сезон новее
			if player.IsNewerSeason(currentSeason, existingDataSeason) {
				// Обновляем физические данные из более свежего турнира
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
					} else {
						logger.Info(ctx, fmt.Sprintf("    🔄 Updated player data: %s (season: %s)", existing.Name, currentSeason))
						updatedCount++
					}
				}
			}

			logger.Info(ctx, fmt.Sprintf("    ♻️  Player EXISTS: %s (ID: %s)", existing.Name, existing.ID))
			existingCount++
			playerIDs = append(playerIDs, existing.ID)
		}
	}

	// Создаем связи player-team-tournament БАТЧЕМ (для новых и существующих игроков)
	if len(playerIDs) > 0 {
		if err := s.CreatePlayerTeamLinksBatch(ctx, playerIDs, teamID, tournamentID, t); err != nil {
			return fmt.Errorf("failed to create player_team links: %w", err)
		}
	}

	logger.Info(ctx, fmt.Sprintf("  📊 ИТОГО: новых=%d, обновлено=%d, существующих=%d, пропущено=%d, слишком старых=%d, всего=%d/%d",
		savedCount, updatedCount, existingCount, skippedCount, skippedTooOld, savedCount+existingCount, len(playersDTO)))
	return nil
}

// MinBirthYear минимальный год рождения игроков для парсинга
// Игроки с годом рождения < MinBirthYear будут пропущены
const MinBirthYear = 2008

// convertPlayerDTO конвертирует DTO в domain entity
func convertPlayerDTO(dto junior.PlayerDTO, season string) (*player.Player, error) {
	// Извлекаем ID из URL
	id := player.ExtractIDFromURL(dto.ProfileURL)
	if id == "" {
		return nil, fmt.Errorf("failed to extract ID from URL: %s", dto.ProfileURL)
	}

	// Парсим дату рождения (формат: 13.05.2008)
	birthDate, err := time.Parse("02.01.2006", dto.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse birth date %s: %w", dto.BirthDate, err)
	}

	// ФИЛЬТРАЦИЯ: пропускаем игроков с годом рождения < MinBirthYear
	if birthDate.Year() < MinBirthYear {
		return nil, fmt.Errorf("birth year %d < %d (too old)", birthDate.Year(), MinBirthYear)
	}

	// Парсим рост
	var height *int
	if dto.Height != "" {
		h, err := strconv.Atoi(strings.TrimSpace(dto.Height))
		if err == nil {
			height = &h
		}
	}

	// Парсим вес
	var weight *int
	if dto.Weight != "" {
		w, err := strconv.Atoi(strings.TrimSpace(dto.Weight))
		if err == nil {
			weight = &w
		}
	}

	// Handedness
	var handedness *string
	if dto.Handedness != "" {
		h := strings.TrimSpace(dto.Handedness)
		handedness = &h
	}

	// DataSeason - сезон из которого взяты данные
	var dataSeason *string
	if season != "" {
		dataSeason = &season
	}

	return &player.Player{
		ID:         id,
		ProfileURL: dto.ProfileURL,
		Name:       strings.TrimSpace(dto.Name),
		BirthDate:  birthDate,
		Position:   strings.TrimSpace(dto.Position),
		Height:     height,
		Weight:     weight,
		Handedness: handedness,
		DataSeason: dataSeason,
		Source:     player.SourceJunior,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}
