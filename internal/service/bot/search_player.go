package bot

import (
	"context"
	"math"

	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// SearchPlayerService сервис поиска игроков для бота
type SearchPlayerService interface {
	Search(ctx context.Context, filters domainBot.SearchFilters, page, pageSize int) (*SearchResult, error)
}

// SearchResult результат поиска с пагинацией
type SearchResult struct {
	Players     []*player.PlayerWithTeam
	TotalCount  int
	TotalPages  int
	CurrentPage int
	PageSize    int
}

// searchPlayerService реализация SearchPlayerService
type searchPlayerService struct {
	playerRepo player.Repository
}

// NewSearchPlayerService создает новый SearchPlayerService
func NewSearchPlayerService(playerRepo player.Repository) SearchPlayerService {
	return &searchPlayerService{
		playerRepo: playerRepo,
	}
}

// Search выполняет поиск игроков по фильтрам бота
func (s *searchPlayerService) Search(ctx context.Context, filters domainBot.SearchFilters, page, pageSize int) (*SearchResult, error) {
	// Конвертируем фильтры бота в фильтры репозитория
	repoFilters := player.SearchFilters{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	// Имя/Фамилия
	if filters.FirstName != nil {
		repoFilters.FirstName = *filters.FirstName
	}
	if filters.LastName != nil {
		repoFilters.LastName = *filters.LastName
	}

	// Год
	if filters.Year != nil {
		repoFilters.BirthYear = filters.Year
	}

	// Позиция
	if filters.Position != nil {
		repoFilters.Position = *filters.Position
	}

	// Рост
	if filters.Height != nil {
		repoFilters.MinHeight = &filters.Height.Min
		repoFilters.MaxHeight = &filters.Height.Max
	}

	// Вес
	if filters.Weight != nil {
		repoFilters.MinWeight = &filters.Weight.Min
		repoFilters.MaxWeight = &filters.Weight.Max
	}

	// Поиск в БД
	players, totalCount, err := s.playerRepo.SearchWithTeam(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	// Рассчитываем общее количество страниц
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
