package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// ExploreHandler handles explore dashboard requests.
type ExploreHandler struct {
	service        *services.ExploreService
	matchesService *services.ExploreMatchesService
}

// NewExploreHandler creates a new explore handler.
func NewExploreHandler(service *services.ExploreService, matchesService *services.ExploreMatchesService) *ExploreHandler {
	return &ExploreHandler{service: service, matchesService: matchesService}
}

// Overview returns platform-wide KPI stats.
func (h *ExploreHandler) Overview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	overview, err := h.service.GetOverview(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to get explore overview: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get overview")
		return
	}
	h.writeJSON(w, http.StatusOK, dto.ExploreOverviewResponse{
		Players:     overview.Players,
		Teams:       overview.Teams,
		Tournaments: overview.Tournaments,
		Matches:     overview.Matches,
	})
}

// Tournaments returns list of tournaments.
func (h *ExploreHandler) Tournaments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	source := r.URL.Query().Get("source")
	domain := r.URL.Query().Get("domain")

	items, err := h.service.GetTournaments(ctx, source, domain)
	if err != nil {
		logger.Error(ctx, "Failed to get tournaments: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get tournaments")
		return
	}

	// Collect IDs and fetch group stats
	ids := make([]string, len(items))
	for i, t := range items {
		ids[i] = t.ID
	}
	groupStats, _ := h.service.GetGroupStats(ctx, ids)

	// Build lookup: tournamentID -> birthYear -> groupName -> stats
	type gsKey struct{ tid, group string; year int }
	gsMap := make(map[gsKey]services.GroupStats, len(groupStats))
	for _, gs := range groupStats {
		gsMap[gsKey{gs.TournamentID, gs.GroupName, gs.BirthYear}] = gs
	}

	tournaments := make([]dto.TournamentItemDTO, len(items))
	for i, t := range items {
		var bygRaw map[string][]string
		if t.BirthYearGroupsRaw != nil {
			_ = json.Unmarshal([]byte(*t.BirthYearGroupsRaw), &bygRaw)
		}

		var byg map[string][]dto.GroupStatsDTO
		if bygRaw != nil {
			byg = make(map[string][]dto.GroupStatsDTO, len(bygRaw))
			for year, groups := range bygRaw {
				yearInt := 0
				if v, err := strconv.Atoi(year); err == nil {
					yearInt = v
				}
				dtos := make([]dto.GroupStatsDTO, len(groups))
				for j, g := range groups {
					gs := gsMap[gsKey{t.ID, g, yearInt}]
					dtos[j] = dto.GroupStatsDTO{Name: g, TeamsCount: gs.TeamsCount, MatchesCount: gs.MatchesCount}
				}
				byg[year] = dtos
			}
		}

		tournaments[i] = dto.TournamentItemDTO{
			ID: t.ID, Name: t.Name, Domain: t.Domain, Season: t.Season,
			Source: t.Source, BirthYearGroups: byg, TeamsCount: t.TeamsCount, MatchesCount: t.MatchesCount, IsEnded: t.IsEnded,
		}
	}
	h.writeJSON(w, http.StatusOK, dto.TournamentListResponse{Tournaments: tournaments})
}

// Standings returns tournament standings.
func (h *ExploreHandler) Standings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournamentID := r.PathValue("id")
	birthYear := parseIntQuery(r, "birthYear", 0)
	groupName := r.URL.Query().Get("group")

	rows, err := h.service.GetTournamentStandings(ctx, tournamentID, birthYear, groupName)
	if err != nil {
		logger.Error(ctx, "Failed to get standings: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get standings")
		return
	}

	standings := make([]dto.StandingDTO, len(rows))
	for i, s := range rows {
		standings[i] = dto.StandingDTO{
			Position: s.Position, Team: s.Team, TeamID: s.TeamID, LogoURL: s.LogoURL,
			Games: s.Games, Wins: s.Wins, WinsOT: s.WinsOT,
			Losses: s.Losses, LossesOT: s.LossesOT, Draws: s.Draws,
			GoalsFor: s.GoalsFor, GoalsAgainst: s.GoalsAgainst, Points: s.Points,
			GroupName: s.GroupName,
		}
	}
	h.writeJSON(w, http.StatusOK, dto.StandingsResponse{Standings: standings})
}

// TournamentMatches returns matches for a tournament.
func (h *ExploreHandler) TournamentMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournamentID := r.PathValue("id")
	limit := parseIntQuery(r, "limit", 30)
	birthYear := parseIntQuery(r, "birthYear", 0)
	groupName := r.URL.Query().Get("group")

	rows, err := h.matchesService.GetTournamentMatches(ctx, tournamentID, birthYear, groupName, limit)
	if err != nil {
		logger.Error(ctx, "Failed to get tournament matches: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get matches")
		return
	}
	h.writeJSON(w, http.StatusOK, dto.MatchListResponse{Matches: matchRowsToDTO(rows)})
}

// Scorers returns top scorers for a tournament.
func (h *ExploreHandler) Scorers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournamentID := r.PathValue("id")
	limit := parseIntQuery(r, "limit", 20)
	birthYear := parseIntQuery(r, "birthYear", 0)
	groupName := r.URL.Query().Get("group")

	rows, err := h.service.GetTournamentScorers(ctx, tournamentID, birthYear, groupName, limit)
	if err != nil {
		logger.Error(ctx, "Failed to get scorers: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get scorers")
		return
	}

	scorers := make([]dto.ScorerDTO, len(rows))
	for i, s := range rows {
		scorers[i] = dto.ScorerDTO{
			Position: i + 1, PlayerID: s.PlayerID, Name: s.Name, PhotoURL: s.PhotoURL,
			Team: s.Team, TeamID: s.TeamID, LogoURL: s.LogoURL,
			Games: s.Games, Goals: s.Goals, Assists: s.Assists, Points: s.Points,
		}
	}
	h.writeJSON(w, http.StatusOK, dto.ScorersResponse{Scorers: scorers})
}

// Seasons returns available seasons list.
func (h *ExploreHandler) Seasons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	seasons, err := h.service.GetSeasons(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to get seasons: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get seasons")
		return
	}
	h.writeJSON(w, http.StatusOK, dto.SeasonsResponse{Seasons: seasons})
}

// TournamentTeams returns teams for a tournament.
func (h *ExploreHandler) TournamentTeams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournamentID := r.PathValue("id")
	birthYear := parseIntQuery(r, "birthYear", 0)
	groupName := r.URL.Query().Get("group")

	rows, err := h.service.GetTournamentTeams(ctx, tournamentID, birthYear, groupName)
	if err != nil {
		logger.Error(ctx, "Failed to get tournament teams: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get teams")
		return
	}

	teams := make([]dto.TeamDTO, len(rows))
	for i, t := range rows {
		teams[i] = dto.TeamDTO{
			ID:           t.ID,
			Name:         t.Name,
			City:         t.City,
			LogoURL:      t.LogoURL,
			PlayersCount: t.PlayersCount,
			GroupName:    t.GroupName,
			BirthYear:    t.BirthYear,
		}
	}
	h.writeJSON(w, http.StatusOK, dto.TeamsResponse{Teams: teams})
}

func (h *ExploreHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *ExploreHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}

// parseIntQuery parses an integer query parameter with a default value.
func parseIntQuery(r *http.Request, key string, defaultVal int) int {
	if s := r.URL.Query().Get(key); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			return v
		}
	}
	return defaultVal
}

// matchRowsToDTO converts service match rows to DTO.
func matchRowsToDTO(rows []services.MatchRow) []dto.MatchDTO {
	matches := make([]dto.MatchDTO, len(rows))
	for i, m := range rows {
		date, timeStr := "", ""
		if m.ScheduledAt != nil {
			date = m.ScheduledAt.Format("2006-01-02")
			timeStr = m.ScheduledAt.Format("15:04")
		}
		matches[i] = dto.MatchDTO{
			ID: m.ID, HomeTeam: m.HomeTeam, AwayTeam: m.AwayTeam,
			HomeTeamID: m.HomeTeamID, AwayTeamID: m.AwayTeamID,
			HomeLogoURL: m.HomeLogoURL, AwayLogoURL: m.AwayLogoURL,
			HomeScore: m.HomeScore, AwayScore: m.AwayScore, ResultType: m.ResultType,
			Date: date, Time: timeStr, Tournament: m.Tournament,
			Venue: m.Venue, Status: m.Status,
		}
	}
	return matches
}
