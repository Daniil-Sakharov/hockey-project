package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/interfaces/http/middleware"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration.
// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "Email and password are required")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "validation_error", "Password must be at least 8 characters")
		return
	}

	user, tokens, err := h.authService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		logger.Error(ctx, "Failed to register user", zap.Error(err))
		if errors.Is(err, services.ErrUserAlreadyExists) {
			writeError(w, http.StatusConflict, "user_exists", "User with this email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to register user")
		return
	}

	writeJSON(w, http.StatusCreated, dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         userToResponse(user),
	})
}

// Login handles user authentication.
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "Email and password are required")
		return
	}

	user, tokens, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.Debug(ctx, "Login failed", zap.Error(err))
		if errors.Is(err, services.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to authenticate")
		return
	}

	writeJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         userToResponse(user),
	})
}

// Refresh handles token refresh.
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "Refresh token is required")
		return
	}

	user, tokens, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		logger.Debug(ctx, "Token refresh failed", zap.Error(err))
		if errors.Is(err, services.ErrInvalidToken) {
			writeError(w, http.StatusUnauthorized, "invalid_token", "Invalid or expired refresh token")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to refresh token")
		return
	}

	writeJSON(w, http.StatusOK, dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         userToResponse(user),
	})
}

// Me returns the current authenticated user.
// GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	user, err := h.authService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		logger.Error(ctx, "Failed to get user", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to get user")
		return
	}

	writeJSON(w, http.StatusOK, userToResponse(user))
}

// LinkPlayer links the authenticated user to a player.
// POST /api/v1/auth/link-player
func (h *AuthHandler) LinkPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	var req dto.LinkPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if req.PlayerID == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "Player ID is required")
		return
	}

	err := h.authService.LinkPlayer(ctx, claims.UserID, req.PlayerID)
	if err != nil {
		logger.Error(ctx, "Failed to link player", zap.Error(err))
		writeError(w, http.StatusBadRequest, "link_failed", err.Error())
		return
	}

	// Return updated user
	user, err := h.authService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to get user")
		return
	}

	writeJSON(w, http.StatusOK, userToResponse(user))
}

// Logout revokes all refresh tokens for the authenticated user.
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims := middleware.GetUserFromContext(ctx)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	err := h.authService.Logout(ctx, claims.UserID)
	if err != nil {
		logger.Error(ctx, "Failed to logout", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to logout")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// userToResponse converts User to UserResponse DTO.
func userToResponse(user *services.User) dto.UserResponse {
	resp := dto.UserResponse{
		ID:                    user.ID,
		Email:                 user.Email,
		PlayerID:              user.PlayerID,
		SubscriptionTier:      user.SubscriptionTier,
		SubscriptionExpiresAt: user.SubscriptionExpiresAt,
		EmailVerified:         user.EmailVerified,
		CreatedAt:             user.CreatedAt,
	}
	if user.Name != nil {
		resp.Name = *user.Name
	}
	return resp
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response.
func writeError(w http.ResponseWriter, status int, errCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error:   errCode,
		Message: message,
	})
}
