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

type tournamentStats struct {
	groups  int
	teams   int
	players int
}

func (o *Orchestrator) processTournamentSafe(ctx context.Context, season dto.SeasonDTO, tournament dto.TournamentDTO) tournamentStats {
	stats, err := o.processTournament(ctx, season, tournament)
	if err != nil {
		logger.Warn(ctx, "Tournament processing failed",
			zap.Int("id", tournament.ID),
			zap.String("name", tournament.Name),
			zap.Error(err),
		)
		return tournamentStats{}
	}
	return stats
}

func (o *Orchestrator) processTournament(ctx context.Context, season dto.SeasonDTO, tournament dto.TournamentDTO) (tournamentStats, error) {
	var stats tournamentStats

	logger.Info(ctx, "Processing tournament",
		zap.Int("id", tournament.ID),
		zap.String("name", tournament.Name),
	)

	// 1. Получаем команды турнира (нужно для запроса групп)
	teams, err := o.fetchTeams(ctx, tournament.ID, 0)
	if err != nil {
		return stats, fmt.Errorf("fetch teams: %w", err)
	}

	if len(teams) == 0 {
		logger.Warn(ctx, "No teams found for tournament",
			zap.Int("tournament_id", tournament.ID),
		)
		return stats, nil
	}

	// 2. Получаем группы турнира через первую команду
	groups, err := o.fetchGroups(ctx, season.ID, tournament.ID, teams[0].ID)
	if err != nil {
		logger.Warn(ctx, "Failed to fetch groups, processing without groups",
			zap.Int("tournament_id", tournament.ID),
			zap.Error(err),
		)
		// Обрабатываем турнир без групп
		groups = []dto.GroupDTO{{ID: 0, Name: ""}}
	}

	// 3. Обрабатываем каждую группу
	for _, group := range groups {
		groupStats := o.processGroup(ctx, season, tournament, group)
		stats.groups++
		stats.teams += groupStats.teams
		stats.players += groupStats.players
	}

	return stats, nil
}

func (o *Orchestrator) fetchTeams(ctx context.Context, tournamentID, groupID int) ([]dto.TeamDTO, error) {
	body := map[string]interface{}{
		"tournament": tournamentID,
	}
	if groupID > 0 {
		body["group"] = groupID
	}

	data, err := o.client.PostAPI("/api/filter/team", body)
	if err != nil {
		return nil, fmt.Errorf("get teams: %w", err)
	}

	teams, err := parsing.ParseTeams(data)
	if err != nil {
		return nil, fmt.Errorf("parse teams: %w", err)
	}

	return teams, nil
}

func (o *Orchestrator) fetchGroups(ctx context.Context, seasonID, tournamentID, teamID int) ([]dto.GroupDTO, error) {
	body := map[string]interface{}{
		"season":     seasonID,
		"tournament": tournamentID,
		"team":       teamID,
	}

	data, err := o.client.PostAPI("/api/filter/data", body)
	if err != nil {
		return nil, fmt.Errorf("get filter data: %w", err)
	}

	filterData, err := parsing.ParseFilterData(data)
	if err != nil {
		return nil, fmt.Errorf("parse filter data: %w", err)
	}

	return filterData.Group, nil
}

type groupStats struct {
	teams   int
	players int
}

func (o *Orchestrator) processGroup(ctx context.Context, season dto.SeasonDTO, tournament dto.TournamentDTO, group dto.GroupDTO) groupStats {
	var stats groupStats

	// Формируем external_id: tournament_group или просто tournament если нет группы
	externalID := strconv.Itoa(tournament.ID)
	if group.ID > 0 {
		externalID = fmt.Sprintf("%d_%d", tournament.ID, group.ID)
	}

	// Сохраняем турнир с группой
	birthYear := tournament.ParseBirthYear()
	seasonStr := tournament.ParseSeason()
	var groupName *string
	if group.Name != "" {
		groupName = &group.Name
	}

	tournamentEntity := &fhmoscowrepo.Tournament{
		ExternalID: externalID,
		Name:       tournament.Name,
		BirthYear:  &birthYear,
		Season:     &seasonStr,
		GroupName:  groupName,
	}

	tournamentDBID, err := o.tournamentRepo.Upsert(ctx, tournamentEntity)
	if err != nil {
		logger.Error(ctx, "Failed to save tournament",
			zap.String("external_id", externalID),
			zap.Error(err),
		)
		return stats
	}

	logger.Info(ctx, "Tournament saved",
		zap.String("id", tournamentDBID),
		zap.String("name", tournament.Name),
		zap.String("group", group.Name),
	)

	// Получаем команды для этой группы
	teams, err := o.fetchTeams(ctx, tournament.ID, group.ID)
	if err != nil {
		logger.Error(ctx, "Failed to fetch teams for group",
			zap.Int("tournament_id", tournament.ID),
			zap.Int("group_id", group.ID),
			zap.Error(err),
		)
		return stats
	}

	// Обрабатываем команды
	stats.teams, stats.players = o.processTeams(ctx, tournamentDBID, teams)

	return stats
}

func (o *Orchestrator) processTeams(ctx context.Context, tournamentID string, teams []dto.TeamDTO) (int, int) {
	var (
		totalTeams   int
		totalPlayers int
		mu           sync.Mutex
	)

	teamCh := make(chan dto.TeamDTO, len(teams))
	for _, t := range teams {
		teamCh <- t
	}
	close(teamCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.TeamWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for team := range teamCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				playersCount := o.processTeamSafe(ctx, tournamentID, team)
				mu.Lock()
				totalTeams++
				totalPlayers += playersCount
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return totalTeams, totalPlayers
}
