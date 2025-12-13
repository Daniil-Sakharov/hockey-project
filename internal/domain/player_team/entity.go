package player_team

import "time"

// PlayerTeam представляет связь игрок-команда-турнир
type PlayerTeam struct {
	PlayerID     string     `db:"player_id"`     // ID игрока
	TeamID       string     `db:"team_id"`       // ID команды
	TournamentID string     `db:"tournament_id"` // ID турнира
	Season       string     `db:"season"`        // Сезон (2025/2026)
	StartedAt    *time.Time `db:"started_at"`    // Когда начал играть
	EndedAt      *time.Time `db:"ended_at"`      // Когда закончил (NULL если активен)
	IsActive     bool       `db:"is_active"`     // Активная связь (игрок сейчас в команде)
	JerseyNumber *int       `db:"jersey_number"` // Номер игрока
	Role         *string    `db:"role"`          // Роль (captain, assistant, player)
	Source       string     `db:"source"`        // Источник данных (junior, registry)
	CreatedAt    time.Time  `db:"created_at"`    // Дата создания
	UpdatedAt    time.Time  `db:"updated_at"`    // Дата обновления
}
