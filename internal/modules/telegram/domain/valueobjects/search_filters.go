package valueobjects

// Range represents a numeric range.
type Range struct {
	Min int
	Max int
}

// SearchFilters contains player search criteria.
type SearchFilters struct {
	FirstName *string
	LastName  *string
	Year      *int
	Position  *string
	Height    *Range
	Weight    *Range
	Region    *string
}

// HasFilters returns true if any filter is set.
func (f *SearchFilters) HasFilters() bool {
	return f.FirstName != nil || f.LastName != nil || f.Year != nil ||
		f.Position != nil || f.Height != nil || f.Weight != nil || f.Region != nil
}

// CountActive returns the number of active filters.
func (f *SearchFilters) CountActive() int {
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

// FIODisplay returns formatted FIO string.
func (f *SearchFilters) FIODisplay() string {
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

// Clear resets all filters.
func (f *SearchFilters) Clear() {
	*f = SearchFilters{}
}
