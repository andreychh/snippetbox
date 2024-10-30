package template

import (
	"time"

	"github.com/andreychh/snippetbox/internal/domain"
)

type Data struct {
	CurrentYear int

	Snippet  domain.Snippet
	Snippets []domain.Snippet
	Form     domain.Form
	Flash    string
}

func NewData(opts ...Option) Data {
	data := Data{
		CurrentYear: time.Now().Year(),
	}
	for _, opt := range opts {
		opt(&data)
	}
	return data
}

type Option func(data *Data)

func WithSnippet(snippet domain.Snippet) Option {
	return func(data *Data) {
		data.Snippet = snippet
	}
}

func WithSnippets(snippets []domain.Snippet) Option {
	return func(data *Data) {
		data.Snippets = snippets
	}
}

func WithForm(form domain.Form) Option {
	return func(data *Data) {
		data.Form = form
	}
}

func WithFlash(flash string) Option {
	return func(data *Data) {
		data.Flash = flash
	}
}
