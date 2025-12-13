package stats

import (
	"fmt"
	"strconv"
	"time"

	clientStats "github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// ConvertToPlayerStatistics конвертирует массив DTO в domain entities
// Теперь валидация и логирование происходят в convertWithTracking
func ConvertToPlayerStatistics(
	dtos []clientStats.PlayerStatisticDTO,
	tournamentID string,
) ([]*player_statistics.PlayerStatistic, error) {
	entities := make([]*player_statistics.PlayerStatistic, 0, len(dtos))

	for _, dto := range dtos {
		entity, err := convertOne(dto, tournamentID)
		if err != nil {
			continue
		}

		// Валидация перед добавлением
		if err := entity.Validate(); err != nil {
			continue
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

// convertOne конвертирует один DTO в entity
func convertOne(dto clientStats.PlayerStatisticDTO, tournamentID string) (*player_statistics.PlayerStatistic, error) {
	// Извлекаем IDs
	playerID := clientStats.ExtractPlayerID(dto.Surname)
	teamID := clientStats.ExtractTeamID(dto.TeamName)

	if playerID == "" || teamID == "" {
		return nil, fmt.Errorf("failed to extract player_id or team_id")
	}

	// Парсим год рождения (с поддержкой заглушек типа "Год_1794")
	birthYear, err := parseBirthYear(dto.BirthYear)
	if err != nil {
		return nil, fmt.Errorf("invalid birth year: %w", err)
	}

	// Создаем entity
	entity := &player_statistics.PlayerStatistic{
		TournamentID: tournamentID,
		PlayerID:     playerID,
		TeamID:       teamID,
		GroupName:    dto.GroupName,
		BirthYear:    birthYear,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),

		// Основная статистика
		Games:          clientStats.ParseInt(dto.GP),
		Goals:          clientStats.ParseInt(dto.G),
		Assists:        clientStats.ParseInt(dto.A),
		Points:         clientStats.ParseInt(dto.PTS),
		Plus:           clientStats.ParseInt(dto.Plus),
		Minus:          clientStats.ParseInt(dto.Minus),
		PlusMinus:      clientStats.ParseInt(dto.PlusMinus),
		PenaltyMinutes: clientStats.ParseInt(dto.PIM),

		// Детальная статистика голов
		GoalsEvenStrength: clientStats.ParseInt(dto.ESG),
		GoalsPowerPlay:    clientStats.ParseInt(dto.PPG),
		GoalsShortHanded:  clientStats.ParseInt(dto.SHG),
		GoalsPeriod1:      clientStats.ParseInt(dto.G1P),
		GoalsPeriod2:      clientStats.ParseInt(dto.G2P),
		GoalsPeriod3:      clientStats.ParseInt(dto.G3P),
		GoalsOvertime:     clientStats.ParseInt(dto.GOT),
		HatTricks:         clientStats.ParseInt(dto.HT),
		GameWinningGoals:  clientStats.ParseInt(dto.WB),

		// Средние показатели
		GoalsPerGame:          clientStats.ParseFloat(dto.GPG),
		PointsPerGame:         clientStats.ParseFloat(dto.PTSPG),
		PenaltyMinutesPerGame: clientStats.ParseFloat(dto.PPM),
	}

	return entity, nil
}

// parseBirthYear парсит год рождения с поддержкой заглушек
// Примеры:
//   - "2009" → 2009
//   - "Год_1794" → 0 (заглушка, невалидный год)
//   - "" → 0 (неизвестный год)
func parseBirthYear(yearStr string) (int, error) {
	if yearStr == "" {
		return 0, nil // Неизвестный год
	}

	// Если начинается с "Год_" - это заглушка для турниров без dropdown годов
	if len(yearStr) > 4 && yearStr[:4] == "Год_" {
		// Возвращаем 0 (неизвестный год) вместо попыток извлечения
		return 0, nil
	}

	// Обычный парсинг года
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, nil // Fallback на 0 вместо ошибки
	}

	// Проверяем что год валидный (2000-2020 для юниорского хоккея)
	if year < 2000 || year > 2020 {
		return 0, nil // Невалидный год → неизвестный
	}

	return year, nil
}
