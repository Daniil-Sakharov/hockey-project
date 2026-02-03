package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// saveMatch сохраняет матч в базу данных
func (o *Orchestrator) saveMatch(ctx context.Context, match dto.MatchDTO,
	season dto.SeasonDTO, tournament dto.TournamentDTO, sub dto.SubTournamentDTO) error {

	tournamentID := fmt.Sprintf("msk:%s-%s-%s", tournament.ID, sub.ID, tournament.GroupID)
	// Team ID format: msk:{tournament_id}-{sub_id}-{group_id}:{team_id}
	homeTeamID := fmt.Sprintf("msk:%s-%s-%s:%s", tournament.ID, sub.ID, tournament.GroupID, match.HomeTeamID)
	awayTeamID := fmt.Sprintf("msk:%s-%s-%s:%s", tournament.ID, sub.ID, tournament.GroupID, match.AwayTeamID)

	status := entities.MatchStatusScheduled
	if match.HomeScore > 0 || match.AwayScore > 0 {
		status = entities.MatchStatusFinished
	}

	entity := &entities.Match{
		ID:           fmt.Sprintf("msk:%s", match.ExternalID),
		ExternalID:   match.ExternalID,
		TournamentID: &tournamentID,
		HomeTeamID:   &homeTeamID,
		AwayTeamID:   &awayTeamID,
		HomeScore:    intPtr(match.HomeScore),
		AwayScore:    intPtr(match.AwayScore),
		MatchNumber:  intPtr(match.MatchNumber),
		ScheduledAt:  timePtr(match.ScheduledAt),
		Status:       status,
		Venue:        strPtr(match.Venue),
		GroupName:    strPtr(sub.Name),
		BirthYear:    intPtr(tournament.BirthYear),
		Source:       Source,
	}

	return o.matchRepo.Upsert(ctx, entity)
}

// saveMatchEvents сохраняет события матча (голы, удаления)
func (o *Orchestrator) saveMatchEvents(ctx context.Context, matchID string, proto *dto.MatchProtocolDTO) (int, error) {
	// Удаляем старые события
	if err := o.matchEventRepo.DeleteByMatchID(ctx, matchID); err != nil {
		return 0, fmt.Errorf("delete old events: %w", err)
	}

	var events []*entities.MatchEvent

	// Конвертируем голы
	for i, goal := range proto.Goals {
		events = append(events, convertGoalToEvent(matchID, goal, i))
	}

	// Конвертируем удаления
	for i, penalty := range proto.Penalties {
		events = append(events, convertPenaltyToEvent(matchID, penalty, i))
	}

	// Сохраняем пакетом
	if len(events) > 0 {
		if err := o.matchEventRepo.CreateBatch(ctx, events); err != nil {
			return 0, fmt.Errorf("create events: %w", err)
		}
	}

	return len(events), nil
}

// saveMatchLineups сохраняет составы матча
func (o *Orchestrator) saveMatchLineups(ctx context.Context, matchID, homeTeamID, awayTeamID, tournamentID string, birthYear int, proto *dto.MatchProtocolDTO) error {
	// Удаляем старые составы
	if err := o.matchLineupRepo.DeleteByMatchID(ctx, matchID); err != nil {
		return fmt.Errorf("delete old lineups: %w", err)
	}

	// Создаём недостающих игроков перед сохранением составов
	o.ensurePlayersExist(ctx, proto.HomeLineup, homeTeamID, tournamentID, birthYear)
	o.ensurePlayersExist(ctx, proto.AwayLineup, awayTeamID, tournamentID, birthYear)

	var lineups []*entities.MatchLineup

	for _, p := range proto.HomeLineup {
		lineups = append(lineups, convertLineupPlayer(matchID, homeTeamID, p))
	}
	for _, p := range proto.AwayLineup {
		lineups = append(lineups, convertLineupPlayer(matchID, awayTeamID, p))
	}

	if len(lineups) > 0 {
		return o.matchLineupRepo.CreateBatch(ctx, lineups)
	}
	return nil
}

// ensurePlayersExist проверяет и создаёт недостающих игроков из протокола
func (o *Orchestrator) ensurePlayersExist(ctx context.Context, lineup []dto.LineupPlayerDTO, teamID, tournamentID string, birthYear int) {
	for _, p := range lineup {
		playerID := fmt.Sprintf("msk:%s", p.PlayerID)

		// Проверяем существует ли игрок
		exists, _ := o.playerRepo.Exists(ctx, p.PlayerID)
		if exists {
			continue
		}

		// Создаём игрока с базовыми данными из протокола
		player := &mihfrepo.Player{
			ExternalID: p.PlayerID,
			FullName:   p.PlayerName,
		}

		// Устанавливаем дату рождения из года турнира
		if birthYear > 0 {
			fallbackDate := time.Date(birthYear, time.January, 1, 0, 0, 0, 0, time.UTC)
			player.BirthDate = &fallbackDate
		}

		// Позиция из протокола
		if p.Position != "" {
			player.Position = &p.Position
		}

		// Сохраняем игрока
		createdID, err := o.playerRepo.Upsert(ctx, player)
		if err != nil {
			logger.Warn(ctx, "Failed to create player from protocol",
				zap.String("player_id", playerID),
				zap.String("name", p.PlayerName),
				zap.Error(err),
			)
			continue
		}

		// Сохраняем связь игрок-команда
		pt := &mihfrepo.PlayerTeam{
			PlayerID:     createdID,
			TeamID:       teamID,
			TournamentID: tournamentID,
		}
		if p.Number > 0 {
			pt.Number = &p.Number
		}
		if p.Position != "" {
			pt.Position = &p.Position
		}

		if err := o.playerTeamRepo.Upsert(ctx, pt); err != nil {
			logger.Warn(ctx, "Failed to create player_team from protocol",
				zap.String("player_id", createdID),
				zap.Error(err),
			)
		}

		logger.Debug(ctx, "Created player from protocol",
			zap.String("player_id", createdID),
			zap.String("name", p.PlayerName),
		)
	}
}
