package providers

import (
	"context"
	"testing"
)

// TestConfig структура для тестирования
type TestConfig struct {
	Name     string `validate:"required"`
	Port     int    `validate:"min=1,max=65535"`
	Enabled  bool
	Optional string
}

// mockConfigSource мок источника конфигурации
type mockConfigSource struct {
	name string
	data map[string]interface{}
	err  error
}

func (m *mockConfigSource) Name() string {
	return m.name
}

func (m *mockConfigSource) Load(ctx context.Context, target interface{}) error {
	if m.err != nil {
		return m.err
	}

	// Простая имитация загрузки данных
	if config, ok := target.(*TestConfig); ok {
		if name, exists := m.data["name"]; exists {
			config.Name = name.(string)
		}
		if port, exists := m.data["port"]; exists {
			config.Port = port.(int)
		}
		if enabled, exists := m.data["enabled"]; exists {
			config.Enabled = enabled.(bool)
		}
	}

	return nil
}

func TestConfigProvider_Load(t *testing.T) {
	tests := []struct {
		name        string
		sources     []ConfigSource
		target      interface{}
		expectError bool
	}{
		{
			name: "successful load with valid config",
			sources: []ConfigSource{
				&mockConfigSource{
					name: "test",
					data: map[string]interface{}{
						"name":    "test-service",
						"port":    8080,
						"enabled": true,
					},
				},
			},
			target:      &TestConfig{},
			expectError: false,
		},
		{
			name: "validation error with invalid config",
			sources: []ConfigSource{
				&mockConfigSource{
					name: "test",
					data: map[string]interface{}{
						"port": 99999, // invalid port
					},
				},
			},
			target:      &TestConfig{},
			expectError: true,
		},
		{
			name:        "error with nil target",
			sources:     []ConfigSource{},
			target:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewConfigProvider(tt.sources...)
			err := provider.Load(context.Background(), tt.target)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Проверяем что данные загрузились корректно
			if !tt.expectError && tt.target != nil {
				config := tt.target.(*TestConfig)
				if config.Name != "test-service" {
					t.Errorf("expected name 'test-service', got '%s'", config.Name)
				}
				if config.Port != 8080 {
					t.Errorf("expected port 8080, got %d", config.Port)
				}
				if !config.Enabled {
					t.Error("expected enabled to be true")
				}
			}
		})
	}
}

func TestConfigProvider_Validate(t *testing.T) {
	provider := NewConfigProvider()

	tests := []struct {
		name        string
		config      interface{}
		expectError bool
	}{
		{
			name: "valid config",
			config: &TestConfig{
				Name: "test",
				Port: 8080,
			},
			expectError: false,
		},
		{
			name: "invalid config - missing required field",
			config: &TestConfig{
				Port: 8080,
				// Name is required but missing
			},
			expectError: true,
		},
		{
			name: "invalid config - port out of range",
			config: &TestConfig{
				Name: "test",
				Port: 99999, // invalid port
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.Validate(context.Background(), tt.config)

			if tt.expectError && err == nil {
				t.Error("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}
