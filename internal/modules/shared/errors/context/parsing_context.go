package context

import (
	"context"
	"time"
)

// ParsingContext контекст для операций парсинга
type ParsingContext struct {
	Source     string                 `json:"source"`      // "junior", "fhspb"
	Domain     string                 `json:"domain"`      // "cfo.fhr.ru", "fhspb.ru"
	EntityType string                 `json:"entity_type"` // "player", "team", "tournament"
	EntityID   string                 `json:"entity_id"`   // внешний ID сущности
	URL        string                 `json:"url"`         // URL который парсится
	UserAgent  string                 `json:"user_agent"`  // User-Agent для запроса
	SessionID  string                 `json:"session_id"`  // ID сессии парсинга
	StartTime  time.Time              `json:"start_time"`  // время начала операции
	Metadata   map[string]interface{} `json:"metadata"`    // дополнительные данные
}

// NewParsingContext создает новый контекст парсинга
func NewParsingContext(source, domain, entityType, entityID, url string) *ParsingContext {
	return &ParsingContext{
		Source:     source,
		Domain:     domain,
		EntityType: entityType,
		EntityID:   entityID,
		URL:        url,
		StartTime:  time.Now(),
		Metadata:   make(map[string]interface{}),
	}
}

// WithUserAgent добавляет User-Agent
func (pc *ParsingContext) WithUserAgent(userAgent string) *ParsingContext {
	pc.UserAgent = userAgent
	return pc
}

// WithSessionID добавляет ID сессии
func (pc *ParsingContext) WithSessionID(sessionID string) *ParsingContext {
	pc.SessionID = sessionID
	return pc
}

// WithMetadata добавляет метаданные
func (pc *ParsingContext) WithMetadata(key string, value interface{}) *ParsingContext {
	pc.Metadata[key] = value
	return pc
}

// Duration возвращает длительность операции
func (pc *ParsingContext) Duration() time.Duration {
	return time.Since(pc.StartTime)
}

// ToMap конвертирует контекст в map для логирования
func (pc *ParsingContext) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"source":      pc.Source,
		"domain":      pc.Domain,
		"entity_type": pc.EntityType,
		"entity_id":   pc.EntityID,
		"url":         pc.URL,
		"start_time":  pc.StartTime,
		"duration_ms": pc.Duration().Milliseconds(),
	}

	if pc.UserAgent != "" {
		result["user_agent"] = pc.UserAgent
	}

	if pc.SessionID != "" {
		result["session_id"] = pc.SessionID
	}

	if len(pc.Metadata) > 0 {
		result["metadata"] = pc.Metadata
	}

	return result
}

// ContextKey тип для ключей контекста
type ContextKey string

const (
	// ParsingContextKey ключ для ParsingContext в context.Context
	ParsingContextKey ContextKey = "parsing_context"
)

// WithParsingContext добавляет ParsingContext в context.Context
func WithParsingContext(ctx context.Context, pc *ParsingContext) context.Context {
	return context.WithValue(ctx, ParsingContextKey, pc)
}

// FromContext извлекает ParsingContext из context.Context
func FromContext(ctx context.Context) (*ParsingContext, bool) {
	pc, ok := ctx.Value(ParsingContextKey).(*ParsingContext)
	return pc, ok
}

// MustFromContext извлекает ParsingContext из context.Context или паникует
func MustFromContext(ctx context.Context) *ParsingContext {
	pc, ok := FromContext(ctx)
	if !ok {
		panic("parsing context not found in context")
	}
	return pc
}
