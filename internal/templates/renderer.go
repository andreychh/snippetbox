package templates

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
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
	var cache = TemplateCache{}
	var pages, err = filepath.Glob("./web/html/pages/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("globbing page templates: %w", err)
	}

	for _, page := range pages {
		var name = filepath.Base(page)
		var templateSet = template.New(name).Funcs(templateFuncs)

		templateSet, err = templateSet.ParseFiles("./web/html/base.gohtml")
		if err != nil {
			return nil, fmt.Errorf("parsing base template: %w", err)
		}

		templateSet, err = templateSet.ParseGlob("./web/html/partials/*.gohtml")
		if err != nil {
			return nil, fmt.Errorf("parsing partial templates: %w", err)
		}

		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("parsing file %q: %w", page, err)
		}

		cache[name] = templateSet
	}

	return cache, nil
}

func (r Renderer) HomePage(data TemplateData) ([]byte, error) {
	return r.renderTemplate("home.gohtml", data)
}

func (r Renderer) SnippetViewPage(data TemplateData) ([]byte, error) {
	return r.renderTemplate("view.gohtml", data)
}

func (r Renderer) SnippetCreatePage(data TemplateData) ([]byte, error) {
	return r.renderTemplate("create.gohtml", data)
}

func (r Renderer) renderTemplate(page string, data TemplateData) ([]byte, error) {
	var templateSet, exists = r.templateCache[page]
	if !exists {
		return nil, ErrNoTemplate
	}

	var buffer bytes.Buffer
	var err = templateSet.ExecuteTemplate(&buffer, "base", data)
	if err != nil {
		return nil, fmt.Errorf("executing template %s: %w", page, err)
	}

	return buffer.Bytes(), nil
}
