package templates

import (
	"html/template"
	"snippetbox/internal/domain"
	"time"
)

type TemplateData struct {
	Snippet  domain.Snippet
	Snippets []domain.Snippet

	CurrentYear int
}

func NewTemplateData() TemplateData {
	return TemplateData{
		CurrentYear: time.Now().Year(),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
