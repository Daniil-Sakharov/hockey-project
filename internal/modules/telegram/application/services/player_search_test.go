package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
)

type mockPlayerRepo struct {
	players []*PlayerWithTeam
	total   int
	err     error
}

func (m *mockPlayerRepo) SearchWithFilters(_ context.Context, _ SearchFilters) ([]*PlayerWithTeam, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.players, m.total, nil
}

func TestPlayerSearchService_Search(t *testing.T) {
	tests := []struct {
		name        string
		players     []*PlayerWithTeam
		total       int
		repoErr     error
		page        int
		pageSize    int
		wantErr     bool
		wantTotal   int
		wantPages   int
		wantCurrent int
	}{
		{
			name: "successful search",
			players: []*PlayerWithTeam{
				{ID: "1", Name: "Иванов Иван"},
				{ID: "2", Name: "Петров Петр"},
			},
			total:       10,
			page:        1,
			pageSize:    5,
			wantErr:     false,
			wantTotal:   10,
			wantPages:   2,
			wantCurrent: 1,
		},
		{
			name:        "empty result",
			players:     []*PlayerWithTeam{},
			total:       0,
			page:        1,
			pageSize:    5,
			wantErr:     false,
			wantTotal:   0,
			wantPages:   1,
			wantCurrent: 1,
		},
		{
			name:    "repository error",
			repoErr: errors.New("db error"),
			page:    1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPlayerRepo{
				players: tt.players,
				total:   tt.total,
				err:     tt.repoErr,
			}
			svc := NewPlayerSearchService(repo)

			result, err := svc.Search(context.Background(), valueobjects.SearchFilters{}, tt.page, tt.pageSize)

			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if result.TotalCount != tt.wantTotal {
				t.Errorf("TotalCount = %d, want %d", result.TotalCount, tt.wantTotal)
			}
			if result.TotalPages != tt.wantPages {
				t.Errorf("TotalPages = %d, want %d", result.TotalPages, tt.wantPages)
			}
			if result.CurrentPage != tt.wantCurrent {
				t.Errorf("CurrentPage = %d, want %d", result.CurrentPage, tt.wantCurrent)
			}
		})
	}
}

func TestPlayerSearchService_FiltersMapping(t *testing.T) {
	repo := &mockPlayerRepo{players: []*PlayerWithTeam{}, total: 0}
	svc := NewPlayerSearchService(repo)

	year := 2010
	position := "Нападающий"
	firstName := "Иван"
	lastName := "Иванов"

	filters := valueobjects.SearchFilters{
		Year:      &year,
		Position:  &position,
		FirstName: &firstName,
		LastName:  &lastName,
	}

	_, err := svc.Search(context.Background(), filters, 1, 10)
	if err != nil {
		t.Errorf("Search() unexpected error: %v", err)
	}
}
