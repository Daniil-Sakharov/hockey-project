package template

// Renderer интерфейс для рендеринга сообщений
type Renderer interface {
	Render(templateName string, data interface{}) (string, error)
}
