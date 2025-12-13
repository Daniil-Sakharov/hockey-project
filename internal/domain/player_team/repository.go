package player_team

import "context"

// Repository интерфейс для работы с player_teams
type Repository interface {
	// Upsert создает или обновляет связь (ON CONFLICT DO UPDATE)
	Upsert(ctx context.Context, pt *PlayerTeam) error

	// UpsertBatch создает или обновляет связи батчем (ON CONFLICT DO UPDATE)
	UpsertBatch(ctx context.Context, pts []*PlayerTeam) error

	// GetByPlayer возвращает все команды игрока
	GetByPlayer(ctx context.Context, playerID string) ([]*PlayerTeam, error)

	// GetActiveByPlayer возвращает активные команды игрока
	GetActiveByPlayer(ctx context.Context, playerID string) ([]*PlayerTeam, error)

	// GetByTeam возвращает всех игроков команды в турнире
	GetByTeam(ctx context.Context, teamID, tournamentID string) ([]*PlayerTeam, error)
}
