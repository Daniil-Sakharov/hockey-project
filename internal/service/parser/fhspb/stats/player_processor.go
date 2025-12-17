package stats

import (
	"context"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

type PlayerProcessor struct {
	deps Dependencies
}

func NewPlayerProcessor(deps Dependencies) *PlayerProcessor {
	return &PlayerProcessor{deps: deps}
}

func (p *PlayerProcessor) Process(ctx context.Context, externalTournamentID int, tournamentID string) error {
	// Загружаем первую страницу
	stats, pageInfo, err := p.deps.Client.GetPlayerStatsFirstPage(ctx, externalTournamentID)
	if err != nil {
		return err
	}

	allStats := stats

	// Если больше одной страницы - загружаем параллельно
	if pageInfo.TotalPages > 1 {
		pageChan := make(chan int, pageInfo.TotalPages-1)
		resultChan := make(chan []dto.PlayerStatsDTO, pageInfo.TotalPages-1)
		errChan := make(chan error, pageInfo.TotalPages-1)

		// Запускаем воркеры
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
					pageStats, err := p.deps.Client.GetPlayerStatsPage(ctx, externalTournamentID, page, pageInfo)
					if err != nil {
						errChan <- err
						return
					}
					resultChan <- pageStats
				}
			}()
		}

		// Отправляем страницы в канал
		for page := 2; page <= pageInfo.TotalPages; page++ {
			pageChan <- page
		}
		close(pageChan)

		wg.Wait()
		close(resultChan)
		close(errChan)

		// Проверяем ошибки
		if err := <-errChan; err != nil {
			return err
		}

		// Собираем результаты
		for pageStats := range resultChan {
			allStats = append(allStats, pageStats...)
		}
	}

	logger.Info(ctx, "loaded player statistics",
		zap.String("tournament_id", tournamentID),
		zap.Int("total_players", len(allStats)),
	)

	// Сохраняем в БД
	return p.saveStats(ctx, tournamentID, allStats)
}

func (p *PlayerProcessor) saveStats(ctx context.Context, tournamentID string, stats []dto.PlayerStatsDTO) error {
	saved := 0
	for _, s := range stats {
		// Получаем player_id по external_id
		player, err := p.deps.PlayerRepo.GetByExternalID(ctx, s.PlayerID)
		if err != nil {
			logger.Debug(ctx, "player not found, skipping", zap.String("player_id", s.PlayerID))
			continue
		}

		// Получаем team_id по external_id и tournament_id
		team, err := p.deps.TeamRepo.GetByExternalID(ctx, s.TeamID, tournamentID)
		if err != nil {
			logger.Debug(ctx, "team not found, skipping", zap.String("team_id", s.TeamID))
			continue
		}

		stat := &fhspbrepo.PlayerStatistics{
			PlayerID:       player.ID,
			TeamID:         team.ID,
			TournamentID:   tournamentID,
			Games:          s.Games,
			Points:         s.Points,
			PointsAvg:      &s.PointsAvg,
			Goals:          s.Goals,
			Assists:        s.Assists,
			PlusMinus:      s.PlusMinus,
			PenaltyMinutes: s.PenaltyMinutes,
			PenaltyAvg:     &s.PenaltyAvg,
		}

		if err := p.deps.PlayerStatisticsRepo.Upsert(ctx, stat); err != nil {
			logger.Error(ctx, "failed to save player statistics", zap.Error(err))
			continue
		}
		saved++
	}

	logger.Info(ctx, "saved player statistics",
		zap.String("tournament_id", tournamentID),
		zap.Int("saved", saved),
		zap.Int("total", len(stats)),
	)

	return nil
}
