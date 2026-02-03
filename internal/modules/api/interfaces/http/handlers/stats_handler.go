package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// StatsHandler handles statistics requests.
type StatsHandler struct {
	statsService *services.StatsService
}

// NewStatsHandler creates a new stats handler.
func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// Overview returns aggregated stats overview.
func (h *StatsHandler) Overview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	overview, err := h.statsService.GetOverview(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to get stats overview: "+err.Error())
		h.writeError(w, http.StatusInternalServerError, "Failed to get statistics")
		return
	}

	response := dto.StatsOverviewResponse{
		Players:     overview.Players,
		Teams:       overview.Teams,
		Tournaments: overview.Tournaments,
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *StatsHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *StatsHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: message})
}
