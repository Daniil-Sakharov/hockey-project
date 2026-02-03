package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenRevoked       = errors.New("token has been revoked")
)

// User represents a user in the system.
type User struct {
	ID                    string     `db:"id"`
	Email                 string     `db:"email"`
	PasswordHash          string     `db:"password_hash"`
	Name                  *string    `db:"name"`
	PlayerID              *string    `db:"player_id"`
	SubscriptionTier      string     `db:"subscription_tier"`
	SubscriptionExpiresAt *time.Time `db:"subscription_expires_at"`
	EmailVerified         bool       `db:"email_verified"`
	LastLoginAt           *time.Time `db:"last_login_at"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
}

// RefreshToken represents a refresh token in the database.
type RefreshToken struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}

// TokenPair represents access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds until access token expires
}

// JWTClaims represents the claims in a JWT token.
type JWTClaims struct {
	UserID           string `json:"user_id"`
	Email            string `json:"email"`
	SubscriptionTier string `json:"subscription_tier"`
	jwt.RegisteredClaims
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// AuthService handles authentication operations.
type AuthService struct {
	db     *sqlx.DB
	config AuthConfig
}

// NewAuthService creates a new auth service.
func NewAuthService(db *sqlx.DB, config AuthConfig) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*User, *TokenPair, error) {
	// Check if user already exists
	var exists bool
	err := s.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		ID:               uuid.New().String(),
		Email:            email,
		PasswordHash:     string(hashedPassword),
		SubscriptionTier: "free",
		EmailVerified:    false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if name != "" {
		user.Name = &name
	}

	query := `
		INSERT INTO users (id, email, password_hash, name, subscription_tier, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = s.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name,
		user.SubscriptionTier, user.EmailVerified, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (*User, *TokenPair, error) {
	// Find user
	var user User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	_, err = s.db.ExecContext(ctx, "UPDATE users SET last_login_at = $1 WHERE id = $2", now, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update last login: %w", err)
	}
	user.LastLoginAt = &now

	// Generate tokens
	tokens, err := s.generateTokens(ctx, &user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &user, tokens, nil
}

// RefreshTokens refreshes the access token using a refresh token.
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*User, *TokenPair, error) {
	// Hash the refresh token to look it up
	tokenHash := hashToken(refreshToken)

	// Find the refresh token
	var rt RefreshToken
	err := s.db.GetContext(ctx, &rt, `
		SELECT * FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`, tokenHash)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	// Find the user
	var user User
	err = s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", rt.UserID)
	if err != nil {
		return nil, nil, ErrUserNotFound
	}

	// Revoke old refresh token
	_, err = s.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1", rt.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate new tokens
	tokens, err := s.generateTokens(ctx, &user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &user, tokens, nil
}

// ValidateAccessToken validates an access token and returns the claims.
func (s *AuthService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUserByID retrieves a user by ID.
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// LinkPlayer links a user account to a player.
func (s *AuthService) LinkPlayer(ctx context.Context, userID, playerID string) error {
	// Verify player exists
	var exists bool
	err := s.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM players WHERE id = $1)", playerID)
	if err != nil {
		return fmt.Errorf("failed to check player existence: %w", err)
	}
	if !exists {
		return errors.New("player not found")
	}

	// Update user
	_, err = s.db.ExecContext(ctx, "UPDATE users SET player_id = $1, updated_at = NOW() WHERE id = $2", playerID, userID)
	if err != nil {
		return fmt.Errorf("failed to link player: %w", err)
	}

	return nil
}

// Logout revokes all refresh tokens for a user.
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL", userID)
	if err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}
	return nil
}

// generateTokens creates a new access and refresh token pair.
func (s *AuthService) generateTokens(ctx context.Context, user *User) (*TokenPair, error) {
	now := time.Now()

	// Create access token
	accessClaims := &JWTClaims{
		UserID:           user.ID,
		Email:            user.Email,
		SubscriptionTier: user.SubscriptionTier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "starrink",
			Subject:   user.ID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token (random UUID)
	refreshTokenString := uuid.New().String()
	refreshTokenHash := hashToken(refreshTokenString)

	// Store refresh token
	rt := &RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: now.Add(s.config.RefreshTokenDuration),
		CreatedAt: now,
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, rt.ID, rt.UserID, rt.TokenHash, rt.ExpiresAt, rt.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.config.AccessTokenDuration.Seconds()),
	}, nil
}

// hashToken creates a SHA256 hash of a token.
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
