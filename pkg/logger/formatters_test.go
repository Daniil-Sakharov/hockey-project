package logger

import "testing"

func TestEnvironmentFormatter(t *testing.T) {
	tests := []struct {
		environment string
		asJSON      bool
	}{
		{"development", false},
		{"staging", false},
		{"production", false},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.environment, func(t *testing.T) {
			encoder := EnvironmentFormatter(tt.environment, tt.asJSON)

			// Проверяем что encoder создался
			if encoder == nil {
				t.Error("Expected encoder, got nil")
			}
		})
	}
}
