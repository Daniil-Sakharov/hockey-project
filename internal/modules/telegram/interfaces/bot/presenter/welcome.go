package presenter

// WelcomeData данные для приветственного сообщения
type WelcomeData struct {
	UserName string
}

// RenderWelcome рендерит приветственное сообщение
func (p *Presenter) RenderWelcome(userName string) (string, error) {
	return p.renderer.Render("welcome.tmpl", WelcomeData{UserName: userName})
}
