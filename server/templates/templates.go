package templates

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/noona-hq/blacklist/logger"
)

func NewRenderer(l logger.Logger) *TemplateRenderer {
	result := &TemplateRenderer{
		logger: l,
	}

	result.templates = result.generateTemplates()

	return result
}

type TemplateRenderer struct {
	logger    logger.Logger
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *TemplateRenderer) generateTemplates() *template.Template {
	success := `server/templates/html/success.html`
	tmpl, err := template.ParseFiles(success)
	if err != nil {
		t.logger.Warnw("Unable to parse templates", "error", err)
	}
	return tmpl
}
