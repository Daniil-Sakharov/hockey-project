package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// CreatePlayerTeamLinksBatch создает связи player-team-tournament батчем
func (s *orchestratorService) CreatePlayerTeamLinksBatch(
	ctx context.Context,
	playerIDs []string,
	teamID string,
	tournamentID string,
	t *entities.Tournament,
) error {
	if len(playerIDs) == 0 {
		return nil
	}

	for _, playerID := range playerIDs {
		pt := &entities.PlayerTeam{
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
		if err := s.playerTeamRepo.Upsert(ctx, pt); err != nil {
			return fmt.Errorf("failed to upsert player_team link: %w", err)
		}
	}

	return nil
}
