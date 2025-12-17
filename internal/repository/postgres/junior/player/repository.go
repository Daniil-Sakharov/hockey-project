package player

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/jmoiron/sqlx"
)

// Проверка реализации интерфейса
var _ player.Repository = (*repository)(nil)

type repository struct {
	db *sqlx.DB
}

// NewRepository создает новый репозиторий игроков
func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}
