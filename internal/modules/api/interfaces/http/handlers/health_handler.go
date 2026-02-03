package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
)

// HealthHandler handles health check requests.
type HealthHandler struct{}

// NewHealthHandler creates a new health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns health status.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := dto.HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}
