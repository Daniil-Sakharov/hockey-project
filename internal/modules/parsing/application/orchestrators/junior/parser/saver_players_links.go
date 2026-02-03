package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// playerLinkContext контекст игрока для создания связи player_team
type playerLinkContext struct {
	JerseyNumber *int
	Height       *int
	Weight       *int
	PhotoURL     string
}

// CreatePlayerTeamLinksBatch создает связи player-team-tournament батчем
func (s *orchestratorService) CreatePlayerTeamLinksBatch(
	ctx context.Context,
	playerIDs []string,
	teamID string,
	tournamentID string,
	t *entities.Tournament,
	birthYear *int,
	groupName *string,
	playerCtx map[string]playerLinkContext,
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
			BirthYear:    birthYear,
			GroupName:    groupName,
			Source:       "junior",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if pctx, ok := playerCtx[playerID]; ok {
			pt.JerseyNumber = pctx.JerseyNumber
			pt.Height = pctx.Height
			pt.Weight = pctx.Weight
			if pctx.PhotoURL != "" {
				pt.PhotoURL = &pctx.PhotoURL
			}
		}

		if err := s.playerTeamRepo.Upsert(ctx, pt); err != nil {
			return fmt.Errorf("failed to upsert player_team link: %w", err)
		}
	}

	return nil
}
