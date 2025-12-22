package presenter

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/template"
)

// Presenter отвечает за форматирование сообщений
type Presenter struct {
	renderer template.Renderer
}

// NewPresenter создает новый Presenter
func NewPresenter(renderer template.Renderer) *Presenter {
	return &Presenter{renderer: renderer}
}
