package detailed

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// Repository интерфейс для сохранения статистики
type Repository interface {
	CreateBatch(ctx context.Context, stats []*entities.PlayerStatistic) (int, error)
}

// SaveOneByOne сохраняет entities по одной с отслеживанием потерь FK
func SaveOneByOne(
	ctx context.Context,
	repo Repository,
	result []*entities.PlayerStatistic,
) (int, []MissingPlayerInfo) {
	savedCount := 0
	losses := []MissingPlayerInfo{}

	for _, entity := range result {
		inserted, err := repo.CreateBatch(ctx, []*entities.PlayerStatistic{entity})
		if err != nil {
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
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
