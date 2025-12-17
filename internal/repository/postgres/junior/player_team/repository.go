package player_team

import (
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

// NewRepository создает новый репозиторий для player_teams
func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}
