package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// SaveTournaments сохраняет турниры в БД (с дедупликацией)
func (s *orchestratorService) SaveTournaments(ctx context.Context, tournamentsDTO []junior.TournamentDTO) ([]*tournament.Tournament, error) {
	var saved []*tournament.Tournament

	for _, dto := range tournamentsDTO {
		// Конвертируем DTO → Entity
		t := convertTournamentDTO(dto)

		// Проверяем существование по URL
		existing, err := s.tournamentRepo.GetByURL(ctx, dto.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing tournament: %w", err)
		}

		if existing != nil {
			// Обновляем существующую запись (могли добавиться новые поля)
			existing.Season = t.Season
			existing.StartDate = t.StartDate
			existing.EndDate = t.EndDate
			existing.IsEnded = t.IsEnded

			if err := s.tournamentRepo.Update(ctx, existing); err != nil {
				return nil, fmt.Errorf("failed to update tournament: %w", err)
			}

			saved = append(saved, existing)
		} else {
			// Создаем новый турнир
			if err := s.tournamentRepo.Create(ctx, t); err != nil {
				return nil, fmt.Errorf("failed to create tournament: %w", err)
			}
			saved = append(saved, t)
		}
	}

	return saved, nil
}

// convertTournamentDTO конвертирует DTO в domain entity
func convertTournamentDTO(dto junior.TournamentDTO) *tournament.Tournament {
	// Парсим дату начала
	var startDate *time.Time
	if dto.StartDate != "" {
		if t, err := time.Parse("02.01.2006", dto.StartDate); err == nil {
			startDate = &t
		}
	}

	// Парсим дату окончания
	var endDate *time.Time
	if dto.EndDate != "" {
		if t, err := time.Parse("02.01.2006", dto.EndDate); err == nil {
			endDate = &t
		}
	}

	return &tournament.Tournament{
		ID:        dto.ID,
		URL:       dto.URL,
		Name:      dto.Name,
		Domain:    dto.Domain,
		Season:    dto.Season,
		StartDate: startDate,
		EndDate:   endDate,
		IsEnded:   dto.IsEnded,
		CreatedAt: time.Now(),
	}
}

// SaveTournamentsBatch сохраняет турниры батчами для оптимизации
func (s *orchestratorService) SaveTournamentsBatch(
	ctx context.Context,
	tournamentsDTO []junior.TournamentDTO,
	batchSize int,
) ([]*tournament.Tournament, error) {
	var allSaved []*tournament.Tournament

	totalCount := len(tournamentsDTO)

	for i := 0; i < totalCount; i += batchSize {
		end := i + batchSize
		if end > totalCount {
			end = totalCount
		}

		batch := tournamentsDTO[i:end]

		// Сохраняем батч
		saved, err := s.SaveTournaments(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to save batch [%d-%d]: %w", i+1, end, err)
		}

		allSaved = append(allSaved, saved...)

		// Логируем прогресс (только если это не последний батч или их больше одного)
		if totalCount > batchSize {
			logger.Info(ctx, fmt.Sprintf("    📦 Saved batch: %d-%d (%d tournaments)", i+1, end, len(saved)))
		}
	}

	return allSaved, nil
}
