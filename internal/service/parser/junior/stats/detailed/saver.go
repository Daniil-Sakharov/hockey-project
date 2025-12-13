package detailed

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// Repository интерфейс для сохранения статистики
type Repository interface {
	CreateBatch(ctx context.Context, stats []*player_statistics.PlayerStatistic) (int, error)
}

// SaveOneByOne сохраняет entities по одной с отслеживанием потерь FK
func SaveOneByOne(
	ctx context.Context,
	repo Repository,
	entities []*player_statistics.PlayerStatistic,
) (int, []MissingPlayerInfo) {
	savedCount := 0
	losses := []MissingPlayerInfo{}

	for _, entity := range entities {
		// Пытаемся сохранить как батч из одной записи
		inserted, err := repo.CreateBatch(ctx, []*player_statistics.PlayerStatistic{entity})
		if err != nil {
			// FK constraint violation - игрок или команда не найдены
			reason := "player_not_found"
			if containsString(err.Error(), "team") || containsString(err.Error(), "teams") {
				reason = "team_not_found"
			}

			losses = append(losses, MissingPlayerInfo{
				PlayerID:  entity.PlayerID,
				TeamID:    entity.TeamID,
				BirthYear: fmt.Sprintf("%d", entity.BirthYear),
				Reason:    reason,
			})
		} else {
			savedCount += inserted
		}
	}

	return savedCount, losses
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
