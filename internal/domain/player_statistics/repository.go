package player_statistics

import "context"

// Repository определяет интерфейс для работы с хранилищем статистики игроков
type Repository interface {
	// CreateBatch создает несколько записей статистики за одну транзакцию
	// Возвращает количество реально вставленных/обновленных записей
	CreateBatch(ctx context.Context, stats []*PlayerStatistic) (int, error)

	// DeleteByTournament удаляет всю статистику турнира
	DeleteByTournament(ctx context.Context, tournamentID string) error

	// DeleteAll удаляет всю статистику (для TRUNCATE)
	DeleteAll(ctx context.Context) error

	// GetByPlayerID возвращает всю статистику игрока
	GetByPlayerID(ctx context.Context, playerID string) ([]*PlayerStatistic, error)

	// GetByTournament возвращает всю статистику турнира
	GetByTournament(ctx context.Context, tournamentID string) ([]*PlayerStatistic, error)

	// CountAll возвращает общее количество записей статистики
	CountAll(ctx context.Context) (int, error)

	// GetAllTimeStats возвращает агрегированную статистику игрока за всё время
	GetAllTimeStats(ctx context.Context, playerID string) (*AggregatedStats, error)

	// GetSeasonStats возвращает агрегированную статистику игрока за сезон
	GetSeasonStats(ctx context.Context, playerID string, season string) (*AggregatedStats, error)

	// GetRecentTournaments возвращает последние N турниров игрока
	// TODO: После исправления парсера будет сортировать по tournaments.start_date
	// Сейчас сортирует по player_statistics.id DESC
	GetRecentTournaments(ctx context.Context, playerID string, limit int) ([]*TournamentStat, error)
}

// AggregatedStats агрегированная статистика (за всё время или за сезон)
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

// TournamentStat статистика по одному турниру с метаданными
type TournamentStat struct {
	TournamentID   string  `db:"tournament_id"`
	TournamentName string  `db:"tournament_name"`
	GroupName      string  `db:"group_name"`
	Season         string  `db:"season"`
	IsChampionship bool    // Заполняется программно, не из БД
	StartDate      *string `db:"start_date"`

	// Статистика
	Games            int `db:"games"`
	Goals            int `db:"goals"`
	Assists          int `db:"assists"`
	Points           int `db:"points"`
	PlusMinus        int `db:"plus_minus"`
	PenaltyMinutes   int `db:"penalty_minutes"`
	HatTricks        int `db:"hat_tricks"`
	GameWinningGoals int `db:"game_winning_goals"`
}
