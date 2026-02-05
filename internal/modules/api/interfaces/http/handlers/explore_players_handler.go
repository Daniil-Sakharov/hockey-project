package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// ExplorePlayersHandler handles player/team explore requests.
type ExplorePlayersHandler struct {
	service *services.ExplorePlayersService
}

// NewExplorePlayersHandler creates a new explore players handler.
func NewExplorePlayersHandler(service *services.ExplorePlayersService) *ExplorePlayersHandler {
	return &ExplorePlayersHandler{service: service}
}

// SearchPlayers handles player search.
func (h *ExplorePlayersHandler) SearchPlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query().Get("q")
	position := r.URL.Query().Get("position")
	season := r.URL.Query().Get("season")
	birthYear := parseIntQuery(r, "birthYear", 0)
	limit := parseIntQuery(r, "limit", 20)
	offset := parseIntQuery(r, "offset", 0)

	rows, total, err := h.service.SearchPlayers(ctx, q, position, season, birthYear, limit, offset)
	if err != nil {
		logger.Error(ctx, "Failed to search players: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to search players")
		return
	}

	players := make([]dto.PlayerItemDTO, len(rows))
	for i, p := range rows {
		players[i] = dto.PlayerItemDTO{
			ID: p.ID, Name: p.Name, Position: p.Position,
			BirthDate: p.BirthDate.Format("2006-01-02"),
			BirthYear: p.BirthDate.Year(),
			Team:      p.Team, TeamID: p.TeamID, TeamLogoURL: p.TeamLogoURL,
			JerseyNumber: p.JerseyNumber, PhotoURL: p.PhotoURL,
			Stats: &dto.PlayerStatsDTO{
				Games: p.Games, Goals: p.Goals, Assists: p.Assists,
				Points: p.Points, PlusMinus: p.PlusMinus, PenaltyMinutes: p.PenaltyMins,
			},
		}
	}
	h.writeJSON(w, http.StatusOK, dto.PlayersSearchResponse{Players: players, Total: total})
}

// PlayerProfile returns a player profile.
func (h *ExplorePlayersHandler) PlayerProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	season := r.URL.Query().Get("season")

	player, err := h.service.GetPlayerProfile(ctx, id, season)
	if err != nil {
		logger.Error(ctx, "Failed to get player profile: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get player")
		return
	}
	if player == nil {
		h.writeError(w, http.StatusNotFound, "Player not found")
		return
	}

	resp := dto.PlayerProfileResponse{
		ID: player.ID, Name: player.Name, Position: player.Position,
		BirthDate: player.BirthDate.Format("2006-01-02"),
		BirthYear: player.BirthDate.Year(),
		Team:      player.Team, TeamID: player.TeamID, TeamLogoURL: player.TeamLogoURL,
		JerseyNumber: player.JerseyNumber,
		Height:       player.Height, Weight: player.Weight,
		Handedness: player.Handedness, City: player.BirthPlace, PhotoURL: player.PhotoURL,
		Stats: &dto.PlayerStatsDTO{
			Games: player.Games, Goals: player.Goals, Assists: player.Assists,
			Points: player.Points, PlusMinus: player.PlusMinus, PenaltyMinutes: player.PenaltyMins,
		},
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// PlayerStats returns detailed stats history for a player.
func (h *ExplorePlayersHandler) PlayerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	rows, err := h.service.GetPlayerStats(ctx, id)
	if err != nil {
		logger.Error(ctx, "Failed to get player stats: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get player stats")
		return
	}

	stats := make([]dto.PlayerStatDTO, len(rows))
	for i, s := range rows {
		stats[i] = dto.PlayerStatDTO{
			Season: s.Season, TournamentID: s.TournamentID, TournamentName: s.TournamentName,
			GroupName: s.GroupName, BirthYear: s.BirthYear,
			Games: s.Games, Goals: s.Goals, Assists: s.Assists,
			Points: s.Points, PlusMinus: s.PlusMinus, PenaltyMinutes: s.PenaltyMinutes,
		}
	}
	h.writeJSON(w, http.StatusOK, dto.PlayerStatsHistoryResponse{Stats: stats})
}

// TeamProfile returns a team profile.
func (h *ExplorePlayersHandler) TeamProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	team, err := h.service.GetTeamProfile(ctx, id)
	if err != nil {
		logger.Error(ctx, "Failed to get team profile: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get team")
		return
	}
	if team == nil {
		h.writeError(w, http.StatusNotFound, "Team not found")
		return
	}

	roster := make([]dto.PlayerItemDTO, len(team.Roster))
	for i, p := range team.Roster {
		roster[i] = dto.PlayerItemDTO{
			ID: p.ID, Name: p.Name, Position: p.Position,
			BirthDate: p.BirthDate.Format("2006-01-02"),
			BirthYear: p.BirthDate.Year(),
			Team:      p.Team, TeamID: p.TeamID, JerseyNumber: p.JerseyNumber,
			PhotoURL: p.PhotoURL,
			Stats: &dto.PlayerStatsDTO{
				Games: p.Games, Goals: p.Goals, Assists: p.Assists,
				Points: p.Points, PlusMinus: p.PlusMinus, PenaltyMinutes: p.PenaltyMins,
			},
		}
	}

	// Convert recent matches
	recentMatches := make([]dto.MatchDTO, len(team.Matches))
	for i, m := range team.Matches {
		date, timeStr := "", ""
		if m.ScheduledAt != nil {
			date = m.ScheduledAt.Format("2006-01-02")
			timeStr = m.ScheduledAt.Format("15:04")
		}
		recentMatches[i] = dto.MatchDTO{
			ID: m.ID, HomeTeam: m.HomeTeam, AwayTeam: m.AwayTeam,
			HomeTeamID: m.HomeTeamID, AwayTeamID: m.AwayTeamID,
			HomeLogoURL: m.HomeLogoURL, AwayLogoURL: m.AwayLogoURL,
			HomeScore: m.HomeScore, AwayScore: m.AwayScore, ResultType: m.ResultType,
			Date: date, Time: timeStr, Tournament: m.Tournament,
			Venue: m.Venue, Status: m.Status,
		}
	}

	resp := dto.TeamProfileResponse{
		ID: team.ID, Name: team.Name, City: team.City, LogoURL: team.LogoURL,
		Tournaments:  team.Tournaments,
		PlayersCount: len(team.Roster),
		Roster:       roster,
		Stats: dto.TeamStatsDTO{
			Wins: team.Stats.Wins, Losses: team.Stats.Losses, Draws: team.Stats.Draws,
			GoalsFor: team.Stats.GoalsFor, GoalsAgainst: team.Stats.GoalsAgainst,
		},
		RecentMatches: recentMatches,
	}
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *ExplorePlayersHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *ExplorePlayersHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}
