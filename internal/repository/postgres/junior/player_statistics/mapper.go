package player_statistics

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// ValidateBeforeInsert проверяет статистику перед вставкой в БД
func ValidateBeforeInsert(stat *player_statistics.PlayerStatistic) error {
	return stat.Validate()
}

// ValidateBatch проверяет массив статистики перед вставкой
func ValidateBatch(stats []*player_statistics.PlayerStatistic) error {
	for _, stat := range stats {
		if err := stat.Validate(); err != nil {
			return err
		}
	}
	return nil
}
