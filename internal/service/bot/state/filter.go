// Package state содержит типы состояния бота
// Реэкспортирует типы из domain/bot для обратной совместимости
package state

import "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"

// Type aliases
type (
	HeightRange   = bot.HeightRange
	WeightRange   = bot.WeightRange
	SearchFilters = bot.SearchFilters
)
