package dto

import "time"

// PlayerProfileDTO профиль игрока с антропометрией
type PlayerProfileDTO struct {
	ID          string     // 19980
	FullName    string     // Заболотный Никита Александрович
	BirthDate   *time.Time // 07.01.2007
	Age         int        // 18
	Height      int        // 200 (см)
	Weight      int        // 98 (кг)
	Position    string     // З (защитник)
	Handedness  string     // Л (левый) / П (правый)
	Citizenship string     // Россия
	PhotoURL    string     // URL фото (если есть)
}
