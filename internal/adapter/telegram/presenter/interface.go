package presenter

// TemplateEngine интерфейс для рендеринга шаблонов
type TemplateEngine interface {
	Render(templateName string, data interface{}) (string, error)
}
