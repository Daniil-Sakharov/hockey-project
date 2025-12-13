package message

import (
	"fmt"
)

// FilterMenuData данные для меню фильтров
type FilterMenuData struct {
	FIO                string // Отображение ФИО (Фамилия Имя или одно из них)
	Year               *int
	Position           *string
	Height             string // форматированный диапазон "150-160"
	Weight             string // форматированный диапазон "50-60"
	Region             *string
	HasFilters         bool
	ActiveFiltersCount int
}

// RenderFilterMenu рендерит меню фильтров
func (p *MessagePresenter) RenderFilterMenu(data FilterMenuData) (string, error) {
	msg, err := p.renderer.Render("filter_menu.tmpl", data)
	if err != nil {
		return "", fmt.Errorf("failed to render filter menu: %w", err)
	}
	return msg, nil
}
