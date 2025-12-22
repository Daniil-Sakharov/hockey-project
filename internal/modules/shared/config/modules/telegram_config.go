package modules

import (
	"fmt"
	"time"
)

// TelegramConfig конфигурация для Telegram бота
type TelegramConfig struct {
	BotToken           string        `env:"TELEGRAM_BOT_TOKEN" validate:"required"`
	WebhookURL         string        `env:"TELEGRAM_WEBHOOK_URL" validate:"omitempty,url"`
	WebhookSecret      string        `env:"TELEGRAM_WEBHOOK_SECRET"`
	Timeout            time.Duration `env:"TELEGRAM_TIMEOUT" default:"30s"`
	UpdateTimeout      int           `env:"TELEGRAM_UPDATE_TIMEOUT" validate:"min=0,max=300" default:"60"`
	MaxConnections     int           `env:"TELEGRAM_MAX_CONNECTIONS" validate:"min=1,max=100" default:"40"`
	AllowedUpdates     []string      `env:"TELEGRAM_ALLOWED_UPDATES"`
	DropPendingUpdates bool          `env:"TELEGRAM_DROP_PENDING_UPDATES" default:"false"`

	// Rate limiting
	RateLimitEnabled bool `env:"TELEGRAM_RATE_LIMIT_ENABLED" default:"true"`
	RateLimitPerSec  int  `env:"TELEGRAM_RATE_LIMIT_PER_SEC" validate:"min=1,max=30" default:"3"`
	RateLimitBurst   int  `env:"TELEGRAM_RATE_LIMIT_BURST" validate:"min=1,max=100" default:"10"`

	// Admin settings
	AdminUserIDs    []int64 `env:"TELEGRAM_ADMIN_USER_IDS"`
	EnableDebugMode bool    `env:"TELEGRAM_DEBUG_MODE" default:"false"`
}

// IsValid проверяет валидность конфигурации Telegram
func (c *TelegramConfig) IsValid() error {
	if c.BotToken == "" {
		return fmt.Errorf("telegram bot token is required")
	}

	// Проверяем формат токена (должен содержать :)
	if len(c.BotToken) < 10 || !contains(c.BotToken, ":") {
		return fmt.Errorf("invalid telegram bot token format")
	}

	return nil
}

// IsWebhookMode проверяет включен ли webhook режим
func (c *TelegramConfig) IsWebhookMode() bool {
	return c.WebhookURL != ""
}

// GetAllowedUpdates возвращает разрешенные типы обновлений
func (c *TelegramConfig) GetAllowedUpdates() []string {
	if len(c.AllowedUpdates) == 0 {
		return []string{"message", "callback_query", "inline_query"}
	}
	return c.AllowedUpdates
}

// contains проверяет содержит ли строка подстроку
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
