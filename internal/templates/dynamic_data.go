package templates

import (
	"html/template"
	"time"

	"github.com/andreychh/snippetbox/internal/domain"
)

type Option func(data *TemplateData)

func WithSnippet(snippet domain.Snippet) Option {
	return func(data *TemplateData) {
		data.Snippet = snippet
	}
}

func WithSnippets(snippets []domain.Snippet) Option {
	return func(data *TemplateData) {
		data.Snippets = snippets
	}
}

type TemplateData struct {
	CurrentYear int

	Snippet  domain.Snippet
	Snippets []domain.Snippet
}

func NewTemplateData(opts ...Option) TemplateData {
	var data = TemplateData{
		CurrentYear: time.Now().Year(),
	}
	for _, opt := range opts {
		opt(&data)
	}
	return data
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
