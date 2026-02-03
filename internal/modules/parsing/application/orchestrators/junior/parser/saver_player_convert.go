package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
)

// MinBirthYear возвращает минимальный год рождения из конфига
func (s *orchestratorService) MinBirthYear() int {
	return s.config.MinBirthYear()
}

// convertPlayerDTO конвертирует DTO в domain entity
func (s *orchestratorService) convertPlayerDTO(dto junior.PlayerDTO, season, domain string) (*entities.Player, error) {
	id := entities.ExtractIDFromURL(dto.ProfileURL)
	if id == "" {
		return nil, fmt.Errorf("failed to extract ID from URL: %s", dto.ProfileURL)
	}

	birthDate, err := time.Parse("02.01.2006", dto.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse birth date %s: %w", dto.BirthDate, err)
	}

	minYear := s.MinBirthYear()
	if birthDate.Year() < minYear {
		return nil, fmt.Errorf("birth year %d < %d (too old)", birthDate.Year(), minYear)
	}

	var height *int
	if dto.Height != "" {
		if h, err := strconv.Atoi(strings.TrimSpace(dto.Height)); err == nil {
			height = &h
		}
	}

	var weight *int
	if dto.Weight != "" {
		if w, err := strconv.Atoi(strings.TrimSpace(dto.Weight)); err == nil {
			weight = &w
		}
	}

	var handedness *string
	if dto.Handedness != "" {
		h := strings.TrimSpace(dto.Handedness)
		handedness = &h
	}

	var citizenship *string
	if dto.Citizenship != "" {
		c := strings.TrimSpace(dto.Citizenship)
		citizenship = &c
	}

	var photoURL *string
	if dto.PhotoURL != "" {
		p := strings.TrimSpace(dto.PhotoURL)
		photoURL = &p
	}

	var dataSeason *string
	if season != "" {
		dataSeason = &season
	}

	var domainPtr *string
	if domain != "" {
		domainPtr = &domain
	}

	return &entities.Player{
		ID:          id,
		ProfileURL:  dto.ProfileURL,
		Name:        strings.TrimSpace(dto.Name),
		BirthDate:   birthDate,
		Position:    strings.TrimSpace(dto.Position),
		Height:      height,
		Weight:      weight,
		Handedness:  handedness,
		Citizenship: citizenship,
		PhotoURL:    photoURL,
		Domain:      domainPtr,
		DataSeason:  dataSeason,
		Source:      entities.SourceJunior,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
