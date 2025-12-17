package bot

import (
	"context"
	"math"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// SearchPlayerService сервис поиска игроков для бота
type SearchPlayerService interface {
	Search(ctx context.Context, filters bot.SearchFilters, page, pageSize int) (*SearchResult, error)
}

// SearchResult результат поиска с пагинацией
type SearchResult struct {
	Players     []*player.PlayerWithTeam
	TotalCount  int
	TotalPages  int
	CurrentPage int
	PageSize    int
}

type searchPlayerService struct {
	playerRepo player.Repository
}

func NewSearchPlayerService(playerRepo player.Repository) SearchPlayerService {
	return &searchPlayerService{playerRepo: playerRepo}
}

func (s *searchPlayerService) Search(ctx context.Context, filters bot.SearchFilters, page, pageSize int) (*SearchResult, error) {
	repoFilters := player.SearchFilters{
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

	players, totalCount, err := s.playerRepo.SearchWithTeam(ctx, repoFilters)
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
