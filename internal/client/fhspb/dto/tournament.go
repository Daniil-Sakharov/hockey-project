package dto

import "time"

// TournamentDTO данные турнира
type TournamentDTO struct {
	ID        int
	Name      string
	BirthYear int
	Season    string
	StartDate *time.Time
	EndDate   *time.Time
	IsEnded   bool
}
