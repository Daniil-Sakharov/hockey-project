package player

import "time"

// Model представляет игрока в БД
type Model struct {
	ID         string    `db:"id"`
	ProfileURL string    `db:"profile_url"`
	Name       string    `db:"name"`
	BirthDate  time.Time `db:"birth_date"`
	Position   string    `db:"position"`
	Height     *int      `db:"height"`
	Weight     *int      `db:"weight"`
	Handedness *string   `db:"handedness"`

	RegistryID *string `db:"registry_id"`
	School     *string `db:"school"`
	Rank       *string `db:"rank"`
	DataSeason *string `db:"data_season"`

	ExternalID  *string `db:"external_id"`
	Citizenship *string `db:"citizenship"`
	Role        *string `db:"role"`
	BirthPlace  *string `db:"birth_place"`

	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ModelWithTeam модель с информацией о команде
type ModelWithTeam struct {
	Model
	TeamName string `db:"team_name"`
	TeamCity string `db:"team_city"`
}
