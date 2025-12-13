package player

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// ToEntity конвертирует DB model в domain entity
func ToEntity(m *Model) *player.Player {
	if m == nil {
		return nil
	}
	return &player.Player{
		ID:          m.ID,
		ProfileURL:  m.ProfileURL,
		Name:        m.Name,
		BirthDate:   m.BirthDate,
		Position:    m.Position,
		Height:      m.Height,
		Weight:      m.Weight,
		Handedness:  m.Handedness,
		RegistryID:  m.RegistryID,
		School:      m.School,
		Rank:        m.Rank,
		DataSeason:  m.DataSeason,
		ExternalID:  m.ExternalID,
		Citizenship: m.Citizenship,
		Role:        m.Role,
		BirthPlace:  m.BirthPlace,
		Source:      m.Source,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToModel конвертирует domain entity в DB model
func ToModel(e *player.Player) *Model {
	if e == nil {
		return nil
	}
	return &Model{
		ID:          e.ID,
		ProfileURL:  e.ProfileURL,
		Name:        e.Name,
		BirthDate:   e.BirthDate,
		Position:    e.Position,
		Height:      e.Height,
		Weight:      e.Weight,
		Handedness:  e.Handedness,
		RegistryID:  e.RegistryID,
		School:      e.School,
		Rank:        e.Rank,
		DataSeason:  e.DataSeason,
		ExternalID:  e.ExternalID,
		Citizenship: e.Citizenship,
		Role:        e.Role,
		BirthPlace:  e.BirthPlace,
		Source:      e.Source,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// ToEntities конвертирует slice моделей
func ToEntities(models []Model) []*player.Player {
	result := make([]*player.Player, len(models))
	for i := range models {
		result[i] = ToEntity(&models[i])
	}
	return result
}

// ToEntityWithTeam конвертирует модель с командой
func ToEntityWithTeam(m *ModelWithTeam) *player.PlayerWithTeam {
	if m == nil {
		return nil
	}
	return &player.PlayerWithTeam{
		Player:   ToEntity(&m.Model),
		TeamName: m.TeamName,
		TeamCity: m.TeamCity,
	}
}
