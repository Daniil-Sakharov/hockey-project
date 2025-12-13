package player_team

import "context"

// Repository определяет интерфейс для работы с хранилищем связей игрок-команда
type Repository interface {
	Upsert(ctx context.Context, pt *PlayerTeam) error
	UpsertBatch(ctx context.Context, pts []*PlayerTeam) error
	GetByPlayer(ctx context.Context, playerID string) ([]*PlayerTeam, error)
	GetByTeam(ctx context.Context, teamID string) ([]*PlayerTeam, error)
	GetByTournament(ctx context.Context, tournamentID string) ([]*PlayerTeam, error)
	GetActiveTeam(ctx context.Context, playerID string) (*PlayerTeam, error)
}
