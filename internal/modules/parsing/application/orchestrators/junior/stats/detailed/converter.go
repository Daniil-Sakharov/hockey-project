package detailed

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
)

// ConvertWithTracking конвертирует DTO в entities с отслеживанием потерь
func ConvertWithTracking(
	dtos []stats.PlayerStatisticDTO,
	tournamentID string,
	convertOne func(stats.PlayerStatisticDTO, string) (*entities.PlayerStatistic, error),
) ([]*entities.PlayerStatistic, []MissingPlayerInfo) {
	result := make([]*entities.PlayerStatistic, 0, len(dtos))
	losses := []MissingPlayerInfo{}

	for _, dto := range dtos {
		entity, err := convertOne(dto, tournamentID)
		if err != nil {
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

		if err := entity.Validate(); err != nil {
			losses = append(losses, MissingPlayerInfo{
				PlayerID:  entity.PlayerID,
				TeamID:    entity.TeamID,
				BirthYear: fmt.Sprintf("%d", entity.BirthYear),
				Reason:    "validation_error",
			})
			continue
		}

		result = append(result, entity)
	}

	return result, losses
}
