package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// RankingHandler handles ranking requests.
type RankingHandler struct {
	rankingService *services.RankingService
}

// NewRankingHandler creates a new ranking handler.
func NewRankingHandler(rankingService *services.RankingService) *RankingHandler {
	return &RankingHandler{rankingService: rankingService}
}

// TopScorers returns top scorers list.
func (h *RankingHandler) TopScorers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse limit from query
	limit := 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	scorers, err := h.rankingService.GetTopScorers(ctx, limit)
	if err != nil {
		logger.Error(ctx, "Failed to get top scorers: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get rankings")
		return
	}

	// Convert to DTO
	players := make([]dto.TopScorerResponse, len(scorers))
	for i, s := range scorers {
		players[i] = dto.TopScorerResponse{
			ID:      s.ID,
			Name:    s.Name,
			Team:    s.Team,
			Goals:   s.Goals,
			Assists: s.Assists,
			Games:   s.Games,
		}
	}

	h.writeJSON(w, http.StatusOK, dto.TopScorersResponse{Players: players})
}

func (h *RankingHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *RankingHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}
