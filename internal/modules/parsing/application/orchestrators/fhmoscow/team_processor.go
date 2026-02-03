package fhmoscow

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	fhmoscowrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhmoscow"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processTeamSafe(ctx context.Context, tournamentID string, team dto.TeamDTO) int {
	count, err := o.processTeam(ctx, tournamentID, team)
	if err != nil {
		logger.Warn(ctx, "Team processing failed",
			zap.Int("id", team.ID),
			zap.String("name", team.Name),
			zap.Error(err),
		)
		return 0
	}
	return count
}

func (o *Orchestrator) processTeam(ctx context.Context, tournamentID string, team dto.TeamDTO) (int, error) {
	logger.Debug(ctx, "Processing team",
		zap.Int("id", team.ID),
		zap.String("name", team.Name),
	)

	// Сохраняем команду
	city := fhmoscowrepo.RegionMoscow
	teamEntity := &fhmoscowrepo.Team{
		ExternalID:   strconv.Itoa(team.ID),
		TournamentID: tournamentID,
		Name:         team.Name,
		City:         &city,
	}

	teamDBID, err := o.teamRepo.Upsert(ctx, teamEntity)
	if err != nil {
		return 0, fmt.Errorf("save team: %w", err)
	}

	// Получаем состав команды
	members, err := o.fetchTeamMembers(ctx, team.ID)
	if err != nil {
		logger.Warn(ctx, "Failed to fetch team members, skipping players",
			zap.Int("team_id", team.ID),
			zap.Error(err),
		)
		return 0, nil
	}

	if len(members) == 0 {
		logger.Debug(ctx, "No members found for team",
			zap.Int("team_id", team.ID),
		)
		return 0, nil
	}

	// Обрабатываем игроков
	playersCount := o.processPlayers(ctx, tournamentID, teamDBID, members)

	logger.Debug(ctx, "Team processed",
		zap.String("team_id", teamDBID),
		zap.Int("players_count", playersCount),
	)

	return playersCount, nil
}

func (o *Orchestrator) fetchTeamMembers(ctx context.Context, teamID int) ([]dto.TeamMemberDTO, error) {
	// Пробуем получить состав через HTML страницу команды
	path := fmt.Sprintf("/team/%d", teamID)
	html, err := o.client.GetHTML(path)
	if err != nil {
		return nil, fmt.Errorf("get team page: %w", err)
	}

	// Сначала пробуем парсить как таблицу
	members, err := parsing.ParseTeamRosterFromTable(html)
	if err != nil || len(members) == 0 {
		// Если не получилось, пробуем извлечь ссылки на игроков
		members, err = parsing.ParseTeamRoster(html)
		if err != nil {
			return nil, fmt.Errorf("parse team roster: %w", err)
		}
	}

	return members, nil
}

func (o *Orchestrator) processPlayers(ctx context.Context, tournamentID, teamID string, members []dto.TeamMemberDTO) int {
	var (
		count int
		mu    sync.Mutex
	)

	memberCh := make(chan dto.TeamMemberDTO, len(members))
	for _, m := range members {
		memberCh <- m
	}
	close(memberCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.PlayerWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for member := range memberCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if o.processPlayerSafe(ctx, tournamentID, teamID, member) {
					mu.Lock()
					count++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return count
}
