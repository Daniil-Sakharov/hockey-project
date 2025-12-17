package bot

// HeightRange диапазон роста
type HeightRange struct {
	Min int
	Max int
}

// WeightRange диапазон веса
type WeightRange struct {
	Min int
	Max int
}

// SearchFilters фильтры поиска игроков
type SearchFilters struct {
	FirstName *string
	LastName  *string
	Year      *int
	Position  *string
	Height    *HeightRange
	Weight    *WeightRange
	Region    *string
}

// HasFilters проверяет есть ли активные фильтры
func (f *SearchFilters) HasFilters() bool {
	return f.FirstName != nil || f.LastName != nil || f.Year != nil ||
		f.Position != nil || f.Height != nil || f.Weight != nil || f.Region != nil
}

// CountActiveFilters возвращает количество активных фильтров
func (f *SearchFilters) CountActiveFilters() int {
	count := 0
	if f.FirstName != nil {
		count++
	}
	if f.LastName != nil {
		count++
	}
	if f.Year != nil {
		count++
	}
	if f.Position != nil {
		count++
	}
	if f.Height != nil {
		count++
	}
	if f.Weight != nil {
		count++
	}
	if f.Region != nil {
		count++
	}
	return count
}

// GetFIODisplay возвращает строку для отображения ФИО
func (f *SearchFilters) GetFIODisplay() string {
	if f.FirstName != nil && f.LastName != nil {
		return *f.LastName + " " + *f.FirstName
	}
	if f.LastName != nil {
		return *f.LastName
	}
	if f.FirstName != nil {
		return *f.FirstName
	}
	return ""
}
