package message

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/template"
)

// MessagePresenter отвечает за рендеринг сообщений через шаблоны
type MessagePresenter struct {
	renderer template.Renderer
}

// NewMessagePresenter создает новый MessagePresenter
func NewMessagePresenter(renderer template.Renderer) *MessagePresenter {
	return &MessagePresenter{
		renderer: renderer,
	}
}
