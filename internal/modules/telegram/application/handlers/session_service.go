package handlers

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
)

// SessionService implements session operations.
type SessionService struct {
	sessions domain.SessionRepository
}

// NewSessionService creates a new session service.
func NewSessionService(sessions domain.SessionRepository) *SessionService {
	return &SessionService{sessions: sessions}
}

// GetOrCreate returns existing session or creates new one.
func (s *SessionService) GetOrCreate(ctx context.Context, userID int64) (*entities.UserSession, error) {
	session, err := s.sessions.Get(userID)
	if err != nil {
		session = entities.NewUserSession(userID)
		if err := s.sessions.Save(session); err != nil {
			return nil, err
		}
	}
	return session, nil
}

// UpdateView updates current view.
func (s *SessionService) UpdateView(ctx context.Context, userID int64, view string) error {
	session, err := s.sessions.Get(userID)
	if err != nil {
		return err
	}
	session.SetView(view)
	return s.sessions.Save(session)
}

// Reset clears session state.
func (s *SessionService) Reset(ctx context.Context, userID int64) error {
	session, err := s.sessions.Get(userID)
	if err != nil {
		return err
	}
	session.Reset()
	return s.sessions.Save(session)
}
