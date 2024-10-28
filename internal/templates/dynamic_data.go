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

func WithForm(form domain.SnippetCreateForm) Option {
	return func(data *TemplateData) {
		data.SnippetCreateForm = form
	}
}

func WithFlash(flash string) Option {
	return func(data *TemplateData) {
		data.Flash = flash
	}
}

type TemplateData struct {
	CurrentYear int

	Snippet           domain.Snippet
	Snippets          []domain.Snippet
	SnippetCreateForm domain.SnippetCreateForm
	Flash             string
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

var templateFuncs = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
