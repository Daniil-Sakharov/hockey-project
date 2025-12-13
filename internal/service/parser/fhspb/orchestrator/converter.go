package orchestrator

import (
	"fmt"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// convertToPlayer конвертирует DTO игрока в domain Player
func convertToPlayer(p *dto.PlayerDTO) (*player.Player, error) {
	birthDate, err := parseDate(p.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("parse birth date: %w", err)
	}

	now := time.Now()

	return &player.Player{
		ID:          p.ExternalID, // Используем ExternalID как первичный ключ
		ProfileURL:  fmt.Sprintf("fhspb://player/%s", p.ExternalID),
		ExternalID:  strPtr(p.ExternalID),
		Source:      player.SourceFHSPB,
		Name:        p.FullName,
		Position:    p.Position,
		Role:        strPtr(p.Role),
		BirthDate:   birthDate,
		BirthPlace:  strPtr(p.BirthPlace),
		Citizenship: strPtr(p.Citizenship),
		Height:      intPtr(p.Height),
		Weight:      intPtr(p.Weight),
		Handedness:  strPtr(p.Stick),
		School:      strPtr(p.School),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse("02.01.2006", s)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}
