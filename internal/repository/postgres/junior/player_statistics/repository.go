package player_statistics

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/jmoiron/sqlx"
)

// Проверка реализации интерфейса
var _ player_statistics.Repository = (*repository)(nil)

type repository struct {
	db *sqlx.DB
}

// NewRepository создает новый репозиторий статистики игроков
func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}
