package services

import (
	"context"
	"math"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
)

// PlayerWithTeam представляет игрока с информацией о команде
type PlayerWithTeam struct {
	ID        string
	Name      string
	BirthDate string
	Position  string
	Height    int
	Weight    int
	TeamName  string
	TeamCity  string
}

// SearchResult результат поиска с пагинацией
type SearchResult struct {
	Players     []*PlayerWithTeam
	TotalCount  int
	TotalPages  int
	CurrentPage int
	PageSize    int
}

// PlayerRepository интерфейс для доступа к данным игроков
type PlayerRepository interface {
	SearchWithFilters(ctx context.Context, filters SearchFilters) ([]*PlayerWithTeam, int, error)
}

// SearchFilters фильтры для поиска в репозитории
type SearchFilters struct {
	FirstName string
	LastName  string
	BirthYear *int
	Position  string
	MinHeight *int
	MaxHeight *int
	MinWeight *int
	MaxWeight *int
	Region    string
	Limit     int
	Offset    int
}

// PlayerSearchService сервис поиска игроков
type PlayerSearchService struct {
	repo PlayerRepository
}

// NewPlayerSearchService создает новый сервис поиска
func NewPlayerSearchService(repo PlayerRepository) *PlayerSearchService {
	return &PlayerSearchService{repo: repo}
}

// Search выполняет поиск игроков по фильтрам
func (s *PlayerSearchService) Search(ctx context.Context, filters valueobjects.SearchFilters, page, pageSize int) (*SearchResult, error) {
	repoFilters := SearchFilters{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	if filters.FirstName != nil {
		repoFilters.FirstName = *filters.FirstName
	}
	if filters.LastName != nil {
		repoFilters.LastName = *filters.LastName
	}
	if filters.Year != nil {
		repoFilters.BirthYear = filters.Year
	}
	if filters.Position != nil {
		repoFilters.Position = *filters.Position
	}
	if filters.Height != nil {
		repoFilters.MinHeight = &filters.Height.Min
		repoFilters.MaxHeight = &filters.Height.Max
	}
	if filters.Weight != nil {
		repoFilters.MinWeight = &filters.Weight.Min
		repoFilters.MaxWeight = &filters.Weight.Max
	}
	if filters.Region != nil {
		repoFilters.Region = *filters.Region
	}

	players, totalCount, err := s.repo.SearchWithFilters(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &SearchResult{
		Players:     players,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}, nil
}
