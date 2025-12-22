package stats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	clientStats "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
)

// ConvertToPlayerStatistics конвертирует массив DTO в domain entities
func ConvertToPlayerStatistics(
	dtos []clientStats.PlayerStatisticDTO,
	tournamentID string,
) ([]*entities.PlayerStatistic, error) {
	result := make([]*entities.PlayerStatistic, 0, len(dtos))

	for _, dto := range dtos {
		entity, err := convertOne(dto, tournamentID)
		if err != nil {
			continue
		}

		if err := entity.Validate(); err != nil {
			continue
		}

		result = append(result, entity)
	}

	return result, nil
}

// convertOne конвертирует один DTO в entity
func convertOne(dto clientStats.PlayerStatisticDTO, tournamentID string) (*entities.PlayerStatistic, error) {
	playerID := clientStats.ExtractPlayerID(dto.Surname)
	teamID := clientStats.ExtractTeamID(dto.TeamName)

	if playerID == "" || teamID == "" {
		return nil, fmt.Errorf("failed to extract player_id or team_id")
	}

	birthYear, _ := parseBirthYear(dto.BirthYear)

	entity := &entities.PlayerStatistic{
		TournamentID: tournamentID,
		PlayerID:     playerID,
		TeamID:       teamID,
		GroupName:    dto.GroupName,
		BirthYear:    birthYear,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),

		Games:          clientStats.ParseInt(dto.GP),
		Goals:          clientStats.ParseInt(dto.G),
		Assists:        clientStats.ParseInt(dto.A),
		Points:         clientStats.ParseInt(dto.PTS),
		Plus:           clientStats.ParseInt(dto.Plus),
		Minus:          clientStats.ParseInt(dto.Minus),
		PlusMinus:      clientStats.ParseInt(dto.PlusMinus),
		PenaltyMinutes: clientStats.ParseInt(dto.PIM),

		GoalsEvenStrength: clientStats.ParseInt(dto.ESG),
		GoalsPowerPlay:    clientStats.ParseInt(dto.PPG),
		GoalsShortHanded:  clientStats.ParseInt(dto.SHG),
		GoalsPeriod1:      clientStats.ParseInt(dto.G1P),
		GoalsPeriod2:      clientStats.ParseInt(dto.G2P),
		GoalsPeriod3:      clientStats.ParseInt(dto.G3P),
		GoalsOvertime:     clientStats.ParseInt(dto.GOT),
		HatTricks:         clientStats.ParseInt(dto.HT),
		GameWinningGoals:  clientStats.ParseInt(dto.WB),

		GoalsPerGame:          clientStats.ParseFloat(dto.GPG),
		PointsPerGame:         clientStats.ParseFloat(dto.PTSPG),
		PenaltyMinutesPerGame: clientStats.ParseFloat(dto.PPM),
	}

	return entity, nil
}

func parseBirthYear(yearStr string) (int, error) {
	if yearStr == "" {
		return 0, nil
	}

	if len(yearStr) > 4 && yearStr[:4] == "Год_" {
		return 0, nil
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, nil
	}

	if year < 2000 || year > 2020 {
		return 0, nil
	}

	return year, nil
}
