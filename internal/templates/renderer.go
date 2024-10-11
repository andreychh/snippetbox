package templates

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var ErrNoTemplate = errors.New("template does not exist")

type TemplateCache map[string]*template.Template

type Renderer struct {
	templateCache TemplateCache
}

func NewRenderer() (Renderer, error) {
	cache, err := newTemplateCache()
	if err != nil {
		return Renderer{}, fmt.Errorf("creating template cache: %w", err)
	}
	return Renderer{templateCache: cache}, nil
}

func newTemplateCache() (TemplateCache, error) {
	var cache = make(TemplateCache)
	var pages, err = filepath.Glob("./web/html/pages/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("globbing page templates: %w", err)
	}

	for _, page := range pages {
		var name = filepath.Base(page)

		var templateSet, err = template.New(name).Funcs(functions).ParseFiles("./web/html/base.gohtml")
		if err != nil {
			return nil, fmt.Errorf("parsing base template: %w", err)
		}

		templateSet, err = templateSet.ParseGlob("./web/html/partials/*.gohtml")
		if err != nil {
			return nil, fmt.Errorf("parsing partial templates: %w", err)
		}

		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("parsing page template %s: %w", name, err)
		}

		cache[name] = templateSet
	}

	return cache, nil
}

func (r Renderer) HomePage(writer http.ResponseWriter, data TemplateData) error {
	return r.render(writer, http.StatusOK, "home.gohtml", data)
}

func (r Renderer) SnippetViewPage(writer http.ResponseWriter, data TemplateData) error {
	return r.render(writer, http.StatusOK, "view.gohtml", data)
}

func (r Renderer) render(writer http.ResponseWriter, status int, page string, data TemplateData) error {
	var templateSet, exists = r.templateCache[page]
	if !exists {
		return ErrNoTemplate
	}

	var buffer bytes.Buffer
	var err = templateSet.ExecuteTemplate(&buffer, "base", data)
	if err != nil {
		return fmt.Errorf("executing template %s: %w", page, err)
	}

	writer.WriteHeader(status)
	_, err = writer.Write(buffer.Bytes())
	return err
}
