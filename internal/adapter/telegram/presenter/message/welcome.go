package message

import (
	"fmt"
)

// WelcomeData данные для приветственного сообщения
type WelcomeData struct {
	UserName         string
	PlayersCount     int
	TeamsCount       int
	TournamentsCount int
}

// RenderWelcome рендерит приветственное сообщение
func (p *MessagePresenter) RenderWelcome(data WelcomeData) (string, error) {
	msg, err := p.renderer.Render("welcome.tmpl", data)
	if err != nil {
		return "", fmt.Errorf("failed to render welcome message: %w", err)
	}
	return msg, nil
}
