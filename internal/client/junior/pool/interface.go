package pool

import (
	jrtypes "github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/types"
)

// SeasonParser интерфейс для парсинга одного сезона
// Реализуется Client'ом, используется Worker Pool'ом
type SeasonParser interface {
	ParseSeasonTournaments(domain, season, ajaxURL string) ([]jrtypes.TournamentDTO, error)
}
