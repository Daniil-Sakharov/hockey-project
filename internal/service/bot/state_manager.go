package bot

import (
	"sync"

	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
)

// stateManager реализация StateManager (in-memory)
type stateManager struct {
	mu     sync.RWMutex
	states map[int64]*domainBot.UserState
}

// NewStateManager создает новый StateManager
func NewStateManager() StateManager {
	return &stateManager{
		states: make(map[int64]*domainBot.UserState),
	}
}

// GetState возвращает состояние пользователя (создает если не существует)
func (sm *stateManager) GetState(userID int64) *domainBot.UserState {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if state, exists := sm.states[userID]; exists {
		return state
	}

	// Создаем новое состояние
	state := &domainBot.UserState{
		UserID:      userID,
		Filters:     domainBot.SearchFilters{},
		CurrentView: "",
		CurrentPage: 1, // По умолчанию первая страница
	}
	sm.states[userID] = state
	return state
}

// UpdateFilters обновляет фильтры пользователя
func (sm *stateManager) UpdateFilters(userID int64, filters domainBot.SearchFilters) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if state, exists := sm.states[userID]; exists {
		state.Filters = filters
	}
}

// SetLastMsgID сохраняет ID последнего сообщения
func (sm *stateManager) SetLastMsgID(userID int64, msgID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if state, exists := sm.states[userID]; exists {
		state.LastMsgID = msgID
	}
}

// SetCurrentView устанавливает текущий экран
func (sm *stateManager) SetCurrentView(userID int64, view string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if state, exists := sm.states[userID]; exists {
		state.CurrentView = view
	}
}

// SetWaitingForInput устанавливает режим ожидания ввода
func (sm *stateManager) SetWaitingForInput(userID int64, input string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if state, exists := sm.states[userID]; exists {
		state.WaitingForInput = input
	}
}

// ResetFilters сбрасывает все фильтры пользователя
func (sm *stateManager) ResetFilters(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if state, exists := sm.states[userID]; exists {
		state.Filters = domainBot.SearchFilters{}
		state.CurrentView = ""
	}
}

// ClearState удаляет состояние пользователя
func (sm *stateManager) ClearState(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.states, userID)
}
