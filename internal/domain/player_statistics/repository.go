package player_statistics

import (
	"context"
	"strings"
)

// Repository определяет интерфейс для работы с хранилищем статистики игроков
type Repository interface {
	CreateBatch(ctx context.Context, stats []*PlayerStatistic) (int, error)
	DeleteByTournament(ctx context.Context, tournamentID string) error
	DeleteAll(ctx context.Context) error
	GetByPlayerID(ctx context.Context, playerID string) ([]*PlayerStatistic, error)
	GetByTournament(ctx context.Context, tournamentID string) ([]*PlayerStatistic, error)
	CountAll(ctx context.Context) (int, error)
	GetAllTimeStats(ctx context.Context, playerID string) (*AggregatedStats, error)
	GetSeasonStats(ctx context.Context, playerID, season string) (*AggregatedStats, error)
	GetRecentTournaments(ctx context.Context, playerID string, limit int) ([]*TournamentStat, error)
}

// AggregatedStats агрегированная статистика
type AggregatedStats struct {
	TournamentsCount int `db:"tournaments_count"`
	Games            int `db:"games"`
	Goals            int `db:"goals"`
	Assists          int `db:"assists"`
	Points           int `db:"points"`
	Plus             int `db:"plus"`
	Minus            int `db:"minus"`
	PlusMinus        int `db:"plus_minus"`
	PenaltyMinutes   int `db:"penalty_minutes"`
	HatTricks        int `db:"hat_tricks"`
	GameWinningGoals int `db:"game_winning_goals"`
}

// TournamentStat статистика по одному турниру
type TournamentStat struct {
	TournamentID     string `db:"tournament_id"`
	TournamentName   string `db:"tournament_name"`
	GroupName        string `db:"group_name"`
	Season           string `db:"season"`
	IsChampionship   bool
	StartDate        *string `db:"start_date"`
	Games            int     `db:"games"`
	Goals            int     `db:"goals"`
	Assists          int     `db:"assists"`
	Points           int     `db:"points"`
	PlusMinus        int     `db:"plus_minus"`
	PenaltyMinutes   int     `db:"penalty_minutes"`
	HatTricks        int     `db:"hat_tricks"`
	GameWinningGoals int     `db:"game_winning_goals"`
}

// CheckIsChampionship проверяет является ли турнир первенством по названию
func (t *TournamentStat) CheckIsChampionship() {
	t.IsChampionship = strings.Contains(strings.ToLower(t.TournamentName), "первенство")
}
