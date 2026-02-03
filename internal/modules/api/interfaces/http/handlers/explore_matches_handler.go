package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// ExploreMatchesHandler handles match and ranking explore requests.
type ExploreMatchesHandler struct {
	service *services.ExploreMatchesService
}

// NewExploreMatchesHandler creates a new explore matches handler.
func NewExploreMatchesHandler(service *services.ExploreMatchesService) *ExploreMatchesHandler {
	return &ExploreMatchesHandler{service: service}
}

// RecentResults returns recently finished matches.
func (h *ExploreMatchesHandler) RecentResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournament := r.URL.Query().Get("tournament")
	limit := parseIntQuery(r, "limit", 20)

	rows, err := h.service.GetRecentResults(ctx, tournament, limit)
	if err != nil {
		logger.Error(ctx, "Failed to get recent results: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get results")
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MatchListResponse{Matches: matchRowsToDTO(rows)})
}

// UpcomingMatches returns upcoming scheduled matches.
func (h *ExploreMatchesHandler) UpcomingMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournament := r.URL.Query().Get("tournament")
	limit := parseIntQuery(r, "limit", 20)

	rows, err := h.service.GetUpcomingMatches(ctx, tournament, limit)
	if err != nil {
		logger.Error(ctx, "Failed to get upcoming matches: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get calendar")
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MatchListResponse{Matches: matchRowsToDTO(rows)})
}

// Rankings returns player rankings sorted by stat.
func (h *ExploreMatchesHandler) Rankings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "points"
	}
	limit := parseIntQuery(r, "limit", 20)
	filter := services.RankingsFilter{
		BirthYear:    parseIntQuery(r, "birthYear", 0),
		Domain:       r.URL.Query().Get("domain"),
		TournamentID: r.URL.Query().Get("tournamentId"),
		GroupName:    r.URL.Query().Get("groupName"),
	}

	result, err := h.service.GetRankings(ctx, sortBy, limit, filter)
	if err != nil {
		logger.Error(ctx, "Failed to get rankings: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get rankings")
		return
	}

	players := make([]dto.RankedPlayerDTO, len(result.Players))
	for i, p := range result.Players {
		players[i] = dto.RankedPlayerDTO{
			Rank: i + 1, ID: p.ID, Name: p.Name, PhotoURL: p.PhotoURL, Position: p.Position,
			BirthYear: p.BirthDate.Year(), Team: p.Team, TeamID: p.TeamID,
			Games: p.Games, Goals: p.Goals, Assists: p.Assists, Points: p.Points,
			PlusMinus: p.PlusMinus, PenaltyMinutes: p.PenaltyMinutes,
		}
	}
	h.writeJSON(w, http.StatusOK, dto.RankingsResponse{Season: result.Season, Players: players})
}

// RankingsFilters returns available filter values for rankings.
func (h *ExploreMatchesHandler) RankingsFilters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.service.GetRankingsFilters(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to get rankings filters: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get filters")
		return
	}
	h.writeJSON(w, http.StatusOK, result)
}

// MatchDetail returns detailed match information.
func (h *ExploreMatchesHandler) MatchDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	result, err := h.service.GetMatchDetail(ctx, id)
	if err != nil {
		logger.Error(ctx, "Failed to get match detail: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get match")
		return
	}
	if result == nil {
		h.writeError(w, http.StatusNotFound, "Match not found")
		return
	}

	m := result.Match
	resp := dto.MatchDetailDTO{
		ID:         m.ID,
		ExternalID: m.ExternalID,
		HomeTeam: dto.MatchTeamDTO{
			ID: m.HomeTeamID, Name: titleCase(m.HomeTeamName),
			City: m.HomeTeamCity, LogoURL: m.HomeLogoURL,
		},
		AwayTeam: dto.MatchTeamDTO{
			ID: m.AwayTeamID, Name: titleCase(m.AwayTeamName),
			City: m.AwayTeamCity, LogoURL: m.AwayLogoURL,
		},
		HomeScore:   m.HomeScore,
		AwayScore:   m.AwayScore,
		ResultType:  m.ResultType,
		Venue:       m.Venue,
		Status:      m.Status,
		GroupName:   m.GroupName,
		BirthYear:   m.BirthYear,
		MatchNumber: m.MatchNumber,
		Tournament: dto.TournamentInfoDTO{
			ID: m.TournamentID, Name: titleCase(m.TournamentName),
		},
	}

	// Date/Time
	if m.ScheduledAt != nil {
		resp.Date = m.ScheduledAt.Format("2006-01-02")
		resp.Time = m.ScheduledAt.Format("15:04")
	}

	// Score by period
	if m.HomeScoreP1 != nil || m.AwayScoreP1 != nil {
		resp.ScoreByPeriod = &dto.ScoreByPeriodDTO{
			HomeP1: m.HomeScoreP1, AwayP1: m.AwayScoreP1,
			HomeP2: m.HomeScoreP2, AwayP2: m.AwayScoreP2,
			HomeP3: m.HomeScoreP3, AwayP3: m.AwayScoreP3,
			HomeOT: m.HomeScoreOT, AwayOT: m.AwayScoreOT,
		}
	}

	// Events
	events := make([]dto.MatchEventDTO, len(result.Events))
	for i, e := range result.Events {
		ev := dto.MatchEventDTO{
			Type:        e.EventType,
			Period:      e.Period,
			TeamName:    titleCase(e.TeamName),
			TeamLogoURL: e.TeamLogoURL,
			GoalType:    e.GoalType,
			PenaltyMins: e.PenaltyMins,
			PenaltyText: e.PenaltyReason,
		}

		// For goals, use scorer info; for penalties, use penalty player info
		if e.EventType == "penalty" {
			ev.PlayerID = e.PenaltyPlayerID
			ev.PlayerName = e.PenaltyPlayerName
			ev.PlayerPhoto = e.PenaltyPlayerPhoto
		} else {
			ev.PlayerID = e.ScorerID
			ev.PlayerName = e.ScorerName
			ev.PlayerPhoto = e.ScorerPhoto
			ev.Assist1ID = e.Assist1ID
			ev.Assist1Name = e.Assist1Name
			ev.Assist2ID = e.Assist2ID
			ev.Assist2Name = e.Assist2Name
		}

		if e.IsHome != nil {
			ev.IsHome = *e.IsHome
		}
		if e.TimeMinutes != nil {
			secs := 0
			if e.TimeSeconds != nil {
				secs = *e.TimeSeconds
			}
			ev.Time = formatGameTime(*e.TimeMinutes, secs)
		}
		events[i] = ev
	}
	resp.Events = events

	// Lineups
	var homeLineup, awayLineup []dto.LineupPlayerDTO
	for _, l := range result.Lineups {
		player := dto.LineupPlayerDTO{
			PlayerID:       l.PlayerID,
			PlayerName:     l.PlayerName,
			PlayerPhoto:    l.PlayerPhoto,
			JerseyNumber:   l.JerseyNumber,
			Position:       mapPositionToAPI(l.Position),
			Goals:          l.Goals,
			Assists:        l.Assists,
			Points:         l.Goals + l.Assists,
			PenaltyMinutes: l.PenaltyMinutes,
			PlusMinus:      l.PlusMinus,
			Saves:          l.Saves,
			GoalsAgainst:   l.GoalsAgainst,
		}
		if l.IsHome {
			homeLineup = append(homeLineup, player)
		} else {
			awayLineup = append(awayLineup, player)
		}
	}
	resp.HomeLineup = homeLineup
	resp.AwayLineup = awayLineup

	h.writeJSON(w, http.StatusOK, resp)
}

func formatGameTime(mins, secs int) string {
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func (h *ExploreMatchesHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *ExploreMatchesHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}
