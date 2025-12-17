package detailed

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// ConvertWithTracking конвертирует DTO в entities с отслеживанием потерь
func ConvertWithTracking(
	dtos []stats.PlayerStatisticDTO,
	tournamentID string,
	convertOne func(stats.PlayerStatisticDTO, string) (*player_statistics.PlayerStatistic, error),
) ([]*player_statistics.PlayerStatistic, []MissingPlayerInfo) {
	entities := make([]*player_statistics.PlayerStatistic, 0, len(dtos))
	losses := []MissingPlayerInfo{}

	for _, dto := range dtos {
		entity, err := convertOne(dto, tournamentID)
		if err != nil {
			// Логируем причину потери
			playerID := stats.ExtractPlayerID(dto.Surname)
			teamID := stats.ExtractTeamID(dto.TeamName)

			losses = append(losses, MissingPlayerInfo{
				PlayerID:  playerID,
				TeamID:    teamID,
				BirthYear: dto.BirthYear,
				Reason:    "conversion_error",
			})
			continue
		}

		// Валидация
		if err := entity.Validate(); err != nil {
			losses = append(losses, MissingPlayerInfo{
				PlayerID:  entity.PlayerID,
				TeamID:    entity.TeamID,
				BirthYear: fmt.Sprintf("%d", entity.BirthYear),
				Reason:    "validation_error",
			})
			continue
		}

		entities = append(entities, entity)
	}

	return entities, losses
}
