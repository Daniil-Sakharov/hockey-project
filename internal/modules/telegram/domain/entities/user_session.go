package entities

import "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"

// UserSession represents a user's session state in the bot.
type UserSession struct {
	UserID          int64
	Filters         valueobjects.SearchFilters
	LastMessageID   int
	CurrentView     string
	WaitingForInput string
	CurrentPage     int
	ResultMsgIDs    []int
	TempFIO         TempFIOData
}

// GetFilters returns pointer to filters
func (s *UserSession) GetFilters() *valueobjects.SearchFilters {
	return &s.Filters
}

// TempFIOData holds temporary FIO input data.
type TempFIOData struct {
	LastName   string
	FirstName  string
	Patronymic string
}

// NewUserSession creates a new user session.
func NewUserSession(userID int64) *UserSession {
	return &UserSession{
		UserID:       userID,
		ResultMsgIDs: make([]int, 0),
	}
}

// Reset clears the session state.
func (s *UserSession) Reset() {
	s.Filters.Clear()
	s.LastMessageID = 0
	s.CurrentView = ""
	s.WaitingForInput = ""
	s.CurrentPage = 0
	s.ResultMsgIDs = nil
	s.TempFIO = TempFIOData{}
}

// IsWaitingForInput returns true if session awaits user input.
func (s *UserSession) IsWaitingForInput() bool {
	return s.WaitingForInput != ""
}

// SetView updates current view and resets page.
func (s *UserSession) SetView(view string) {
	s.CurrentView = view
	s.CurrentPage = 0
}

// NextPage increments current page.
func (s *UserSession) NextPage() {
	s.CurrentPage++
}

// PrevPage decrements current page if possible.
func (s *UserSession) PrevPage() {
	if s.CurrentPage > 0 {
		s.CurrentPage--
	}
}

// ApplyTempFIO applies temporary FIO data to filters.
func (s *UserSession) ApplyTempFIO() {
	// Копируем значения, а не указатели на поля TempFIO
	if s.TempFIO.LastName != "" {
		lastName := s.TempFIO.LastName
		s.Filters.LastName = &lastName
	} else {
		s.Filters.LastName = nil
	}
	if s.TempFIO.FirstName != "" {
		firstName := s.TempFIO.FirstName
		s.Filters.FirstName = &firstName
	} else {
		s.Filters.FirstName = nil
	}
	s.TempFIO = TempFIOData{}
}

// ClearTempFIO clears temporary FIO data.
func (s *UserSession) ClearTempFIO() {
	s.TempFIO = TempFIOData{}
}

// ResetFilters clears all filters.
func (s *UserSession) ResetFilters() {
	s.Filters.Clear()
	s.CurrentPage = 0
}
