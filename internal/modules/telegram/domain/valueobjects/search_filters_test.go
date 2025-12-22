package valueobjects

import "testing"

func TestSearchFilters_HasFilters(t *testing.T) {
	tests := []struct {
		name    string
		filters SearchFilters
		want    bool
	}{
		{"empty", SearchFilters{}, false},
		{"with year", SearchFilters{Year: ptr(2010)}, true},
		{"with position", SearchFilters{Position: strPtr("Нападающий")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filters.HasFilters(); got != tt.want {
				t.Errorf("HasFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchFilters_CountActive(t *testing.T) {
	tests := []struct {
		name    string
		filters SearchFilters
		want    int
	}{
		{"empty", SearchFilters{}, 0},
		{"one filter", SearchFilters{Year: ptr(2010)}, 1},
		{"two filters", SearchFilters{Year: ptr(2010), Position: strPtr("Защитник")}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filters.CountActive(); got != tt.want {
				t.Errorf("CountActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchFilters_FIODisplay(t *testing.T) {
	tests := []struct {
		name    string
		filters SearchFilters
		want    string
	}{
		{"empty", SearchFilters{}, ""},
		{"lastname only", SearchFilters{LastName: strPtr("Иванов")}, "Иванов"},
		{"firstname only", SearchFilters{FirstName: strPtr("Петр")}, "Петр"},
		{"both", SearchFilters{LastName: strPtr("Иванов"), FirstName: strPtr("Петр")}, "Иванов Петр"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filters.FIODisplay(); got != tt.want {
				t.Errorf("FIODisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptr(i int) *int          { return &i }
func strPtr(s string) *string { return &s }
