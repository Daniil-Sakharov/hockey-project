package dto

import "time"

// PlayerProfileDTO профиль игрока со страницы /player/{id}
type PlayerProfileDTO struct {
	ID         string     // external_id
	FullName   string     // ФИО
	BirthDate  *time.Time // дата рождения
	Age        int        // возраст
	Position   string     // позиция (В, З, Н)
	Height     int        // рост в см
	Weight     int        // вес в кг
	Handedness string     // хват (Л, П)
	PhotoURL   string     // URL фото

	// Статистика по сезонам
	Stats []PlayerStatsDTO
}

// TeamMemberDTO член команды из страницы /team/{id}
type TeamMemberDTO struct {
	PlayerID string // ID игрока
	Name     string // имя
	Number   int    // номер
	Position string // позиция
}
