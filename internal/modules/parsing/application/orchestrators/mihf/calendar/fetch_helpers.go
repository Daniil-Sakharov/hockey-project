package calendar

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// fetchGroups получает группы сезона
func (o *Orchestrator) fetchGroups(ctx context.Context, season dto.SeasonDTO) ([]dto.GroupDTO, error) {
	html, err := o.client.Get(season.URL)
	if err != nil {
		return nil, fmt.Errorf("get season page: %w", err)
	}

	groups, err := parsing.ParseGroups(html, season.Year)
	if err != nil {
		return nil, fmt.Errorf("parse groups: %w", err)
	}

	logger.Debug(ctx, "Groups fetched",
		zap.String("season", season.Year),
		zap.Int("count", len(groups)),
	)
	return groups, nil
}

// fetchTournaments получает турниры группы
func (o *Orchestrator) fetchTournaments(ctx context.Context, season dto.SeasonDTO, group dto.GroupDTO) ([]dto.TournamentDTO, error) {
	html, err := o.client.Get(group.URL)
	if err != nil {
		return nil, fmt.Errorf("get group page: %w", err)
	}

	tournaments, err := parsing.ParseTournaments(html, season.Year, group.ID)
	if err != nil {
		return nil, fmt.Errorf("parse tournaments: %w", err)
	}

	logger.Debug(ctx, "Tournaments fetched",
		zap.String("group", group.Name),
		zap.Int("count", len(tournaments)),
	)
	return tournaments, nil
}

// fetchSubTournaments получает подтурниры турнира
func (o *Orchestrator) fetchSubTournaments(ctx context.Context, tournament dto.TournamentDTO) ([]dto.SubTournamentDTO, error) {
	html, err := o.client.Get(tournament.URL)
	if err != nil {
		return nil, fmt.Errorf("get tournament page: %w", err)
	}

	subTournaments, err := parsing.ParseSubTournaments(html, "", tournament.GroupID, tournament.ID)
	if err != nil {
		return nil, fmt.Errorf("parse sub-tournaments: %w", err)
	}

	logger.Debug(ctx, "Sub-tournaments fetched",
		zap.String("tournament", tournament.Name),
		zap.Int("count", len(subTournaments)),
	)
	return subTournaments, nil
}
