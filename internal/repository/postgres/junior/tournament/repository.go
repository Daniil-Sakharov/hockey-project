package tournament

import (
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

// NewRepository создает новый репозиторий турниров
func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}
