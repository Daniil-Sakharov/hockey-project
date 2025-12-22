package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// SaveTournaments сохраняет турниры в БД
func (s *orchestratorService) SaveTournaments(ctx context.Context, tournamentsDTO []junior.TournamentDTO) ([]*entities.Tournament, error) {
	var saved []*entities.Tournament

	for _, dto := range tournamentsDTO {
		t := convertTournamentDTO(dto)

		existing, err := s.tournamentRepo.GetByURL(ctx, dto.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing tournament: %w", err)
		}

		if existing != nil {
			existing.Season = t.Season
			existing.StartDate = t.StartDate
			existing.EndDate = t.EndDate
			existing.IsEnded = t.IsEnded

			if err := s.tournamentRepo.Update(ctx, existing); err != nil {
				return nil, fmt.Errorf("failed to update tournament: %w", err)
			}
			saved = append(saved, existing)
		} else {
			if err := s.tournamentRepo.Create(ctx, t); err != nil {
				return nil, fmt.Errorf("failed to create tournament: %w", err)
			}
			saved = append(saved, t)
		}
	}

	return saved, nil
}

func convertTournamentDTO(dto junior.TournamentDTO) *entities.Tournament {
	var startDate *time.Time
	if dto.StartDate != "" {
		if t, err := time.Parse("02.01.2006", dto.StartDate); err == nil {
			startDate = &t
		}
	}

	var endDate *time.Time
	if dto.EndDate != "" {
		if t, err := time.Parse("02.01.2006", dto.EndDate); err == nil {
			endDate = &t
		}
	}

	return &entities.Tournament{
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

// SaveTournamentsBatch сохраняет турниры батчами
func (s *orchestratorService) SaveTournamentsBatch(
	ctx context.Context,
	tournamentsDTO []junior.TournamentDTO,
	batchSize int,
) ([]*entities.Tournament, error) {
	var allSaved []*entities.Tournament

	totalCount := len(tournamentsDTO)

	for i := 0; i < totalCount; i += batchSize {
		end := i + batchSize
		if end > totalCount {
			end = totalCount
		}

		saved, err := s.SaveTournaments(ctx, tournamentsDTO[i:end])
		if err != nil {
			return nil, fmt.Errorf("failed to save batch [%d-%d]: %w", i+1, end, err)
		}

		allSaved = append(allSaved, saved...)

		if totalCount > batchSize {
			logger.Info(ctx, fmt.Sprintf("    📦 Saved batch: %d-%d", i+1, end))
		}
	}

	return allSaved, nil
}
