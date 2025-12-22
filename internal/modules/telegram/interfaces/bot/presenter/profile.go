package presenter

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
)

// RenderProfile рендерит профиль игрока
func (p *Presenter) RenderProfile(profile *services.PlayerProfile) (string, error) {
	return p.renderer.Render("player_profile.tmpl", profile)
}
