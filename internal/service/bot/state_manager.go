package bot

import (
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
)

type stateManager struct {
	mu     sync.RWMutex
	states map[int64]*bot.UserState
}

func NewStateManager() StateManager {
	return &stateManager{
		states: make(map[int64]*bot.UserState),
	}
}

func (sm *stateManager) GetState(userID int64) *bot.UserState {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, exists := sm.states[userID]; exists {
		return s
	}

	s := &bot.UserState{
		UserID:      userID,
		Filters:     bot.SearchFilters{},
		CurrentView: "",
		CurrentPage: 1,
	}
	sm.states[userID] = s
	return s
}

func (sm *stateManager) UpdateFilters(userID int64, filters bot.SearchFilters) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, exists := sm.states[userID]; exists {
		s.Filters = filters
	}
}

func (sm *stateManager) SetLastMsgID(userID int64, msgID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, exists := sm.states[userID]; exists {
		s.LastMsgID = msgID
	}
}

func (sm *stateManager) SetCurrentView(userID int64, view string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, exists := sm.states[userID]; exists {
		s.CurrentView = view
	}
}

func (sm *stateManager) SetWaitingForInput(userID int64, input string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, exists := sm.states[userID]; exists {
		s.WaitingForInput = input
	}
}

func (sm *stateManager) ResetFilters(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, exists := sm.states[userID]; exists {
		s.Filters = bot.SearchFilters{}
		s.CurrentView = ""
	}
}

func (sm *stateManager) ClearState(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.states, userID)
}
