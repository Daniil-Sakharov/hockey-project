package template

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/telegram/*.tmpl
var templateFS embed.FS

// Engine рендерер шаблонов на основе text/template
type Engine struct {
	templates *template.Template
}

// NewEngine создает новый Engine
func NewEngine() (*Engine, error) {
	// Создаем template с функциями
	funcMap := template.FuncMap{
		"pluralizeGoals":        pluralizeGoals,
		"pluralizeWinningGoals": pluralizeWinningGoals,
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/telegram/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Engine{
		templates: tmpl,
	}, nil
}

// Render рендерит шаблон с данными
func (e *Engine) Render(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := e.templates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}
	return buf.String(), nil
}

// pluralizeGoals склоняет слово "гол"
func pluralizeGoals(count int) string {
	if count%10 == 1 && count%100 != 11 {
		return "гол"
	}
	if (count%10 >= 2 && count%10 <= 4) && (count%100 < 10 || count%100 >= 20) {
		return "гола"
	}
	return "голов"
}

// pluralizeWinningGoals склоняет "победный гол"
func pluralizeWinningGoals(count int) string {
	if count%10 == 1 && count%100 != 11 {
		return "победный гол"
	}
	if (count%10 >= 2 && count%10 <= 4) && (count%100 < 10 || count%100 >= 20) {
		return "победных гола"
	}
	return "победных голов"
}
