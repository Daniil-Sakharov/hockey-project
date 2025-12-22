package presenter

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
)

// FilterMenuData данные для меню фильтров
type FilterMenuData struct {
	FIO                string
	Year               *int
	Position           *string
	Height             string
	Weight             string
	Region             *string
	HasFilters         bool
	ActiveFiltersCount int
}

// RenderFilterMenu рендерит меню фильтров
func (p *Presenter) RenderFilterMenu(filters *valueobjects.SearchFilters) (string, error) {
	data := FilterMenuData{
		FIO:                filters.FIODisplay(),
		Year:               filters.Year,
		Position:           filters.Position,
		Region:             filters.Region,
		HasFilters:         filters.HasFilters(),
		ActiveFiltersCount: filters.CountActive(),
	}

	if filters.Height != nil {
		data.Height = fmt.Sprintf("%d-%d", filters.Height.Min, filters.Height.Max)
	}
	if filters.Weight != nil {
		data.Weight = fmt.Sprintf("%d-%d", filters.Weight.Min, filters.Weight.Max)
	}

	return p.renderer.Render("filter_menu.tmpl", data)
}

// RenderFilterMenuFromSession рендерит меню фильтров из сессии (для возврата)
func (p *Presenter) RenderFilterMenuFromSession(session interface {
	GetFilters() *valueobjects.SearchFilters
},
) string {
	filters := session.GetFilters()
	text, err := p.RenderFilterMenu(filters)
	if err != nil {
		return "🔍 *Поиск игроков*\n\nВыберите параметры поиска:"
	}
	return text
}
