package fhmoscow

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	fhmoscowrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhmoscow"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

type scanStats struct {
	scanned   int64
	saved     int64
	skipped   int64
	notFound  int64
	errors    int64
}

// scanPlayers сканирует игроков по ID и сохраняет их профили
func (o *Orchestrator) scanPlayers(ctx context.Context) scanStats {
	var stats scanStats
	maxID := o.config.MaxPlayerID()

	logger.Info(ctx, "[STEP 4] Scanning players by ID...",
		zap.Int("max_player_id", maxID),
		zap.Int("min_birth_year", o.config.MinBirthYear()),
		zap.Int("workers", o.config.PlayerWorkers()),
	)

	// Канал для ID игроков
	idCh := make(chan int, 1000)
	go func() {
		for id := 1; id <= maxID; id++ {
			select {
			case <-ctx.Done():
				close(idCh)
				return
			case idCh <- id:
			}
		}
		close(idCh)
	}()

	var wg sync.WaitGroup
	for i := 0; i < o.config.PlayerWorkers(); i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			o.scanPlayerWorker(ctx, workerID, idCh, &stats)
		}(i + 1)
	}

	wg.Wait()

	logger.Info(ctx, "[STEP 4] Player scan completed",
		zap.Int64("scanned", stats.scanned),
		zap.Int64("saved", stats.saved),
		zap.Int64("skipped", stats.skipped),
		zap.Int64("not_found", stats.notFound),
		zap.Int64("errors", stats.errors),
	)

	return stats
}

func (o *Orchestrator) scanPlayerWorker(ctx context.Context, workerID int, idCh <-chan int, stats *scanStats) {
	for playerID := range idCh {
		select {
		case <-ctx.Done():
			return
		default:
		}

		atomic.AddInt64(&stats.scanned, 1)

		profile, err := o.fetchPlayerProfileByID(ctx, playerID)
		if err != nil {
			// 404 = игрок не существует
			if isNotFoundError(err) {
				atomic.AddInt64(&stats.notFound, 1)
				continue
			}
			atomic.AddInt64(&stats.errors, 1)
			logger.Debug(ctx, "Failed to fetch player profile",
				zap.Int("player_id", playerID),
				zap.Error(err),
			)
			continue
		}

		// Фильтруем по году рождения
		if profile.BirthDate != nil {
			birthYear := profile.BirthDate.Year()
			if birthYear < o.config.MinBirthYear() {
				atomic.AddInt64(&stats.skipped, 1)
				continue
			}
		}

		// Сохраняем игрока
		if err := o.saveScannedPlayer(ctx, profile); err != nil {
			atomic.AddInt64(&stats.errors, 1)
			logger.Warn(ctx, "Failed to save player",
				zap.Int("player_id", playerID),
				zap.String("name", profile.FullName),
				zap.Error(err),
			)
			continue
		}

		atomic.AddInt64(&stats.saved, 1)

		// Логируем прогресс каждые 100 игроков
		scanned := atomic.LoadInt64(&stats.scanned)
		if scanned%500 == 0 {
			logger.Info(ctx, "Scan progress",
				zap.Int64("scanned", scanned),
				zap.Int64("saved", atomic.LoadInt64(&stats.saved)),
				zap.Int("worker_id", workerID),
			)
		}
	}
}

func (o *Orchestrator) fetchPlayerProfileByID(ctx context.Context, playerID int) (*dto.PlayerProfileDTO, error) {
	path := fmt.Sprintf("/player/%d", playerID)
	html, err := o.client.GetHTML(path)
	if err != nil {
		return nil, fmt.Errorf("get player page: %w", err)
	}

	profile, err := parsing.ParsePlayerProfile(html, fmt.Sprintf("%d", playerID))
	if err != nil {
		return nil, fmt.Errorf("parse profile: %w", err)
	}

	return profile, nil
}

func (o *Orchestrator) saveScannedPlayer(ctx context.Context, profile *dto.PlayerProfileDTO) error {
	profileURL := fmt.Sprintf("/player/%s", profile.ID)

	player := &fhmoscowrepo.Player{
		ExternalID: profile.ID,
		FullName:   profile.FullName,
		ProfileURL: &profileURL,
		BirthDate:  profile.BirthDate,
	}

	if profile.Position != "" {
		player.Position = &profile.Position
	}
	if profile.Height > 0 {
		player.Height = &profile.Height
	}
	if profile.Weight > 0 {
		player.Weight = &profile.Weight
	}
	if profile.Handedness != "" {
		player.Handedness = &profile.Handedness
	}

	// 1. Сохраняем игрока
	playerID, err := o.playerRepo.Upsert(ctx, player)
	if err != nil {
		return fmt.Errorf("save player: %w", err)
	}

	// 2. Сохраняем статистику для каждой записи
	for _, stat := range profile.Stats {
		if err := o.savePlayerStat(ctx, playerID, profile.Position, stat); err != nil {
			logger.Debug(ctx, "Failed to save player stat",
				zap.String("player_id", playerID),
				zap.String("team", stat.TeamName),
				zap.String("tournament", stat.TournamentName),
				zap.Error(err),
			)
			// Продолжаем с другими записями
			continue
		}
	}

	return nil
}

func (o *Orchestrator) savePlayerStat(ctx context.Context, playerID, position string, stat dto.PlayerStatsDTO) error {
	var teamID, tournamentID string
	var err error

	// 1. Пробуем найти существующую команду по external_id и названию турнира
	if stat.TeamID > 0 {
		existingTeam, findErr := o.teamRepo.FindByExternalIDAndTournamentName(
			ctx,
			fmt.Sprintf("%d", stat.TeamID),
			stat.TournamentName,
		)
		if findErr == nil && existingTeam != nil {
			teamID = existingTeam.ID
			tournamentID = existingTeam.TournamentID
		}
	}

	// 2. Если не нашли по external_id — пробуем по имени команды
	if teamID == "" && stat.TeamName != "" {
		existingTeam, findErr := o.teamRepo.FindByNameAndTournamentName(
			ctx,
			stat.TeamName,
			stat.TournamentName,
		)
		if findErr == nil && existingTeam != nil {
			teamID = existingTeam.ID
			tournamentID = existingTeam.TournamentID
		}
	}

	// 3. Если не нашли — fallback на создание турнира и команды
	if tournamentID == "" {
		tournamentID, err = o.findOrCreateTournament(ctx, stat.TournamentName, stat.Season)
		if err != nil {
			return fmt.Errorf("find/create tournament: %w", err)
		}
	}

	if teamID == "" {
		teamID, err = o.findOrCreateTeam(ctx, stat.TeamID, stat.TeamName, tournamentID)
		if err != nil {
			return fmt.Errorf("find/create team: %w", err)
		}
	}

	// 3. Сохраняем связь player_team
	pt := &fhmoscowrepo.PlayerTeam{
		PlayerID:     playerID,
		TeamID:       teamID,
		TournamentID: tournamentID,
		Season:       &stat.Season,
	}
	if position != "" {
		pt.Position = &position
	}

	if err := o.playerTeamRepo.Upsert(ctx, pt); err != nil {
		return fmt.Errorf("save player_team: %w", err)
	}

	// 4. Сохраняем статистику
	ps := &fhmoscowrepo.PlayerStatistics{
		PlayerID:       playerID,
		TeamID:         teamID,
		TournamentID:   tournamentID,
		Games:          stat.Games,
		Goals:          stat.Goals,
		Assists:        stat.Assists,
		Points:         stat.Points,
		PenaltyMinutes: stat.PenaltyMinutes,
	}

	if err := o.playerStatisticsRepo.Upsert(ctx, ps); err != nil {
		return fmt.Errorf("save player_statistics: %w", err)
	}

	return nil
}

func (o *Orchestrator) findOrCreateTournament(ctx context.Context, name, season string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty tournament name")
	}

	// Извлекаем год рождения из названия турнира (например "ПМ 2009 г.р. 25/26" -> 2009)
	birthYear := extractBirthYearFromName(name)

	tournament := &fhmoscowrepo.Tournament{
		ExternalID: fmt.Sprintf("scanned:%s:%s", season, name),
		Name:       name,
		Season:     &season,
		BirthYear:  &birthYear,
	}

	return o.tournamentRepo.Upsert(ctx, tournament)
}

func (o *Orchestrator) findOrCreateTeam(ctx context.Context, teamID int, teamName, tournamentID string) (string, error) {
	if teamName == "" {
		return "", fmt.Errorf("empty team name")
	}

	externalID := fmt.Sprintf("%d", teamID)
	if teamID == 0 {
		// Если нет ID, используем имя
		externalID = fmt.Sprintf("name:%s", teamName)
	}

	team := &fhmoscowrepo.Team{
		ExternalID:   externalID,
		TournamentID: tournamentID,
		Name:         teamName,
	}

	return o.teamRepo.Upsert(ctx, team)
}

func extractBirthYearFromName(name string) int {
	// Ищем 4-значное число в названии (год рождения)
	// Например "ПМ 2009 г.р. 25/26" -> 2009
	for i := 0; i <= len(name)-4; i++ {
		if name[i] >= '0' && name[i] <= '9' &&
			name[i+1] >= '0' && name[i+1] <= '9' &&
			name[i+2] >= '0' && name[i+2] <= '9' &&
			name[i+3] >= '0' && name[i+3] <= '9' {
			year := int(name[i]-'0')*1000 + int(name[i+1]-'0')*100 + int(name[i+2]-'0')*10 + int(name[i+3]-'0')
			if year >= 2005 && year <= 2020 {
				return year
			}
		}
	}
	return 0
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "404") || contains(errStr, "not found")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsImpl(s, substr))
}

func containsImpl(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
