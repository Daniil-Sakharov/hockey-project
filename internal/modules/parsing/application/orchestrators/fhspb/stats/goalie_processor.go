package stats

import (
	"context"
	"sync"

	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

type GoalieProcessor struct {
	deps Dependencies
}

func NewGoalieProcessor(deps Dependencies) *GoalieProcessor {
	return &GoalieProcessor{deps: deps}
}

func (p *GoalieProcessor) Process(ctx context.Context, externalTournamentID int, tournamentID string) error {
	stats, pageInfo, err := p.deps.Client.GetGoalieStatsFirstPage(ctx, externalTournamentID)
	if err != nil {
		return err
	}

	allStats := stats

	if pageInfo.TotalPages > 1 {
		pageChan := make(chan int, pageInfo.TotalPages-1)
		resultChan := make(chan []dto.GoalieStatsDTO, pageInfo.TotalPages-1)
		errChan := make(chan error, pageInfo.TotalPages-1)

		var wg sync.WaitGroup
		workers := p.deps.StatisticsWorkers
		if workers <= 0 {
			workers = 3
		}

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for page := range pageChan {
					pageStats, err := p.deps.Client.GetGoalieStatsPage(ctx, externalTournamentID, page, pageInfo)
					if err != nil {
						errChan <- err
						return
					}
					resultChan <- pageStats
				}
			}()
		}

		for page := 2; page <= pageInfo.TotalPages; page++ {
			pageChan <- page
		}
		close(pageChan)

		wg.Wait()
		close(resultChan)
		close(errChan)

		if err := <-errChan; err != nil {
			return err
		}

		for pageStats := range resultChan {
			allStats = append(allStats, pageStats...)
		}
	}

	logger.Info(ctx, "loaded goalie statistics",
		zap.String("tournament_id", tournamentID),
		zap.Int("total_goalies", len(allStats)),
	)

	return p.saveStats(ctx, tournamentID, allStats)
}

func (p *GoalieProcessor) saveStats(ctx context.Context, tournamentID string, stats []dto.GoalieStatsDTO) error {
	saved := 0
	for _, s := range stats {
		player, err := p.deps.PlayerRepo.GetByExternalID(ctx, s.PlayerID)
		if err != nil {
			logger.Debug(ctx, "goalie not found, skipping", zap.String("player_id", s.PlayerID))
			continue
		}

		team, err := p.deps.TeamRepo.GetByExternalID(ctx, s.TeamID, tournamentID)
		if err != nil {
			logger.Debug(ctx, "team not found, skipping", zap.String("team_id", s.TeamID))
			continue
		}

		stat := &fhspbrepo.GoalieStatistics{
			PlayerID:        player.ID,
			TeamID:          team.ID,
			TournamentID:    tournamentID,
			Games:           s.Games,
			Minutes:         s.Minutes,
			GoalsAgainst:    s.GoalsAgainst,
			ShotsAgainst:    s.ShotsAgainst,
			SavePercentage:  &s.SavePercentage,
			GoalsAgainstAvg: &s.GoalsAgainstAvg,
			Wins:            s.Wins,
			Shutouts:        s.Shutouts,
			Assists:         s.Assists,
			PenaltyMinutes:  s.PenaltyMinutes,
		}

		if err := p.deps.GoalieStatisticsRepo.Upsert(ctx, stat); err != nil {
			logger.Error(ctx, "failed to save goalie statistics", zap.Error(err))
			continue
		}
		saved++
	}

	logger.Info(ctx, "saved goalie statistics",
		zap.String("tournament_id", tournamentID),
		zap.Int("saved", saved),
		zap.Int("total", len(stats)),
	)

	return nil
}
