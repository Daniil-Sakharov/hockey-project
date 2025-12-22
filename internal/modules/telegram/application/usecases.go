package application

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
)

// FilterUseCase handles filter operations.
type FilterUseCase interface {
	ApplyFilter(ctx context.Context, userID int64, filterType string, value interface{}) error
	ResetFilters(ctx context.Context, userID int64) error
	GetFilters(ctx context.Context, userID int64) (*valueobjects.SearchFilters, error)
}

// SessionUseCase handles session operations.
type SessionUseCase interface {
	GetOrCreate(ctx context.Context, userID int64) (*entities.UserSession, error)
	UpdateView(ctx context.Context, userID int64, view string) error
	Reset(ctx context.Context, userID int64) error
}

// SearchUseCase handles player search operations.
type SearchUseCase interface {
	SearchPlayers(ctx context.Context, userID int64, page int) (SearchResult, error)
}

// SearchResult contains search results with pagination.
type SearchResult struct {
	Players    []PlayerDTO
	TotalCount int
	Page       int
	HasMore    bool
}

// PlayerDTO represents player data for presentation.
type PlayerDTO struct {
	ID       int64
	Name     string
	Year     int
	Position string
	Height   int
	Weight   int
	Region   string
}
