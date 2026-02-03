package dto

import "time"

// RegisterRequest represents user registration request.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

// LoginRequest represents user login request.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest represents token refresh request.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents authentication response with tokens.
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"` // seconds
	User         UserResponse `json:"user"`
}

// UserResponse represents user data in responses.
type UserResponse struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	Name             string     `json:"name"`
	PlayerID         *string    `json:"player_id,omitempty"`
	SubscriptionTier string     `json:"subscription_tier"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	EmailVerified    bool       `json:"email_verified"`
	CreatedAt        time.Time  `json:"created_at"`
}

// LinkPlayerRequest represents request to link user to player.
type LinkPlayerRequest struct {
	PlayerID string `json:"player_id" validate:"required"`
}

// UpdateProfileRequest represents profile update request.
type UpdateProfileRequest struct {
	Name string `json:"name" validate:"omitempty,min=2"`
}
