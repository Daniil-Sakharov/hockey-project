package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// CreatePlayerTeamLinksBatch создает связи player-team-tournament батчем
func (s *orchestratorService) CreatePlayerTeamLinksBatch(
	ctx context.Context,
	playerIDs []string,
	teamID string,
	tournamentID string,
	t *tournament.Tournament,
) error {
	if len(playerIDs) == 0 {
		return nil
	}

	// Формируем батч
	links := make([]*player_team.PlayerTeam, 0, len(playerIDs))

	for _, playerID := range playerIDs {
		pt := &player_team.PlayerTeam{
			PlayerID:     playerID,
			TeamID:       teamID,
			TournamentID: tournamentID,
			Season:       t.Season,
			StartedAt:    t.StartDate,
			EndedAt:      t.EndDate,
			IsActive:     !t.IsEnded,
			Source:       "junior",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		links = append(links, pt)
	}

	// Batch INSERT
	if err := s.playerTeamRepo.UpsertBatch(ctx, links); err != nil {
		return fmt.Errorf("failed to batch upsert player_team links: %w", err)
	}

	return nil
}
