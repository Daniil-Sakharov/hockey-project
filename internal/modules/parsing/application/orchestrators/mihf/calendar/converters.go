package calendar

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
)

// extractCityFromTeamName извлекает город из названия команды
// Формат: "Локомотив-2004 (Ярославль)" → "Ярославль"
func extractCityFromTeamName(teamName string) string {
	re := regexp.MustCompile(`\(([^)]+)\)\s*$`)
	matches := re.FindStringSubmatch(teamName)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// convertGoalToEvent конвертирует гол в событие матча
func convertGoalToEvent(matchID string, goal dto.GoalEventDTO, idx int) *entities.MatchEvent {
	eventID := fmt.Sprintf("%s:goal:%d", matchID, idx)
	return &entities.MatchEvent{
		ID:              eventID,
		MatchID:         matchID,
		EventType:       "goal",
		Period:          intPtr(goal.Period),
		TimeMinutes:     intPtr(goal.TimeMinutes),
		TimeSeconds:     intPtr(goal.TimeSeconds),
		ScorerPlayerID:  strPtr(fmt.Sprintf("msk:%s", goal.ScorerID)),
		Assist1PlayerID: strPtrIfNotEmpty(goal.Assist1ID),
		Assist2PlayerID: strPtrIfNotEmpty(goal.Assist2ID),
		GoalType:        strPtr(goal.GoalType),
		ScoreHome:       intPtr(goal.ScoreAfter[0]),
		ScoreAway:       intPtr(goal.ScoreAfter[1]),
		IsHome:          &goal.IsHome,
		Source:          Source,
	}
}

// convertPenaltyToEvent конвертирует удаление в событие матча
func convertPenaltyToEvent(matchID string, penalty dto.PenaltyEventDTO, idx int) *entities.MatchEvent {
	eventID := fmt.Sprintf("%s:penalty:%d", matchID, idx)
	return &entities.MatchEvent{
		ID:              eventID,
		MatchID:         matchID,
		EventType:       "penalty",
		Period:          intPtr(penalty.Period),
		TimeMinutes:     intPtr(penalty.TimeMinutes),
		TimeSeconds:     intPtr(penalty.TimeSeconds),
		PenaltyPlayerID: strPtr(fmt.Sprintf("msk:%s", penalty.PlayerID)),
		PenaltyMinutes:  intPtr(penalty.Minutes),
		PenaltyReason:   strPtr(penalty.Reason),
		IsHome:          &penalty.IsHome,
		Source:          Source,
	}
}

// convertLineupPlayer конвертирует игрока состава
func convertLineupPlayer(matchID, teamID string, p dto.LineupPlayerDTO) *entities.MatchLineup {
	playerID := fmt.Sprintf("msk:%s", p.PlayerID)
	return &entities.MatchLineup{
		ID:           fmt.Sprintf("%s:%s", matchID, playerID),
		MatchID:      matchID,
		PlayerID:     playerID,
		TeamID:       teamID,
		JerseyNumber: intPtr(p.Number),
		Position:     strPtr(p.Position),
		CaptainRole:  strPtr(p.CaptainRole),
		Source:       Source,
	}
}

// Helper functions
func intPtr(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

// scorePtr возвращает указатель на int, включая 0 (для счёта)
func scorePtr(v int) *int {
	return &v
}

func strPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func strPtrIfNotEmpty(id string) *string {
	if id == "" {
		return nil
	}
	s := fmt.Sprintf("msk:%s", id)
	return &s
}

func timePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
