package services

import (
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
)

// UserStateService управляет состоянием пользователей
type UserStateService struct {
	states map[int64]*entities.UserSession
	mu     sync.RWMutex
}

// NewUserStateService создает новый сервис состояния
func NewUserStateService() *UserStateService {
	return &UserStateService{
		states: make(map[int64]*entities.UserSession),
	}
}

// GetSession возвращает сессию пользователя (создает если не существует)
func (s *UserStateService) GetSession(userID int64) *entities.UserSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, ok := s.states[userID]; ok {
		return session
	}

	session := entities.NewUserSession(userID)
	s.states[userID] = session
	return session
}

// GetFilters возвращает фильтры пользователя
func (s *UserStateService) GetFilters(userID int64) *valueobjects.SearchFilters {
	session := s.GetSession(userID)
	return &session.Filters
}

// UpdateFilters обновляет фильтры пользователя
func (s *UserStateService) UpdateFilters(userID int64, filters *valueobjects.SearchFilters) {
	session := s.GetSession(userID)
	if filters != nil {
		session.Filters = *filters
	}
}

// GetCurrentPage возвращает текущую страницу
func (s *UserStateService) GetCurrentPage(userID int64) int {
	session := s.GetSession(userID)
	if session.CurrentPage < 1 {
		return 1
	}
	return session.CurrentPage
}

// SetCurrentPage устанавливает текущую страницу
func (s *UserStateService) SetCurrentPage(userID int64, page int) {
	session := s.GetSession(userID)
	session.CurrentPage = page
}

// SetCurrentView устанавливает текущий view
func (s *UserStateService) SetCurrentView(userID int64, view string) {
	session := s.GetSession(userID)
	session.CurrentView = view
}

// SetLastMessageID устанавливает ID последнего сообщения
func (s *UserStateService) SetLastMessageID(userID int64, msgID int) {
	session := s.GetSession(userID)
	session.LastMessageID = msgID
}

// SetWaitingForInput устанавливает ожидание ввода
func (s *UserStateService) SetWaitingForInput(userID int64, input string) {
	session := s.GetSession(userID)
	session.WaitingForInput = input
}

// ResetFilters сбрасывает фильтры
func (s *UserStateService) ResetFilters(userID int64) {
	session := s.GetSession(userID)
	session.ResetFilters()
}

// ClearSession очищает сессию пользователя
func (s *UserStateService) ClearSession(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, userID)
}
