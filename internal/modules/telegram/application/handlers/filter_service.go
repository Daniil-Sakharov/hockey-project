package handlers

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
)

// FilterService implements filter operations.
type FilterService struct {
	sessions domain.SessionRepository
}

// NewFilterService creates a new filter service.
func NewFilterService(sessions domain.SessionRepository) *FilterService {
	return &FilterService{sessions: sessions}
}

// ApplyFilter applies a filter to user session.
func (s *FilterService) ApplyFilter(ctx context.Context, userID int64, filterType string, value interface{}) error {
	session, err := s.sessions.Get(userID)
	if err != nil {
		return err
	}

	switch filterType {
	case "year":
		if v, ok := value.(int); ok {
			session.Filters.Year = &v
		}
	case "position":
		if v, ok := value.(string); ok {
			session.Filters.Position = &v
		}
	case "region":
		if v, ok := value.(string); ok {
			session.Filters.Region = &v
		}
	case "height":
		if v, ok := value.(valueobjects.Range); ok {
			session.Filters.Height = &v
		}
	case "weight":
		if v, ok := value.(valueobjects.Range); ok {
			session.Filters.Weight = &v
		}
	case "lastname":
		if v, ok := value.(string); ok {
			session.Filters.LastName = &v
		}
	case "firstname":
		if v, ok := value.(string); ok {
			session.Filters.FirstName = &v
		}
	}

	return s.sessions.Save(session)
}

// ResetFilters clears all filters.
func (s *FilterService) ResetFilters(ctx context.Context, userID int64) error {
	session, err := s.sessions.Get(userID)
	if err != nil {
		return err
	}
	session.Filters.Clear()
	return s.sessions.Save(session)
}

// GetFilters returns current filters.
func (s *FilterService) GetFilters(ctx context.Context, userID int64) (*valueobjects.SearchFilters, error) {
	session, err := s.sessions.Get(userID)
	if err != nil {
		return nil, err
	}
	return &session.Filters, nil
}
