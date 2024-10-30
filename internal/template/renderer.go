package template

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"
)

type Page string

const (
	PageHome     Page = "home.gohtml"
	PageView     Page = "view.gohtml"
	PageCreate   Page = "create.gohtml"
	PageNotFound Page = "not_found.gohtml"
)

type Renderer struct {
	cache Cache
}

func NewRenderer() (*Renderer, error) {
	cache, err := newCache()
	if err != nil {
		return nil, fmt.Errorf("creating template cache: %w", err)
	}
	return &Renderer{cache: cache}, nil
}

func (r Renderer) RenderPage(page Page, data Data) ([]byte, error) {
	templateSet, exists := r.cache[string(page)]
	if !exists {
		return nil, fmt.Errorf("template %q is not found in cache", page)
	}

	buffer := new(bytes.Buffer)
	err := templateSet.ExecuteTemplate(buffer, "base", data)
	if err != nil {
		return nil, fmt.Errorf("executing template %q: %w", page, err)
	}

	return buffer.Bytes(), nil
}

type Cache map[string]*template.Template

func newCache() (Cache, error) {
	cache := Cache{}
	files, err := filepath.Glob("./web/html/pages/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("globbing '.gohtml' files: %w", err)
	}

	for _, file := range files {
		name := filepath.Base(file)
		templateSet := template.New(name).Funcs(templateFuncs)

		templateSet, err = templateSet.ParseFiles("./web/html/base.gohtml")
		if err != nil {
			return nil, fmt.Errorf("parsing base template: %w", err)
		}

		templateSet, err = templateSet.ParseGlob("./web/html/partials/*.gohtml")
		if err != nil {
			return nil, fmt.Errorf("parsing partial template: %w", err)
		}

		templateSet, err = templateSet.ParseFiles(file)
		if err != nil {
			return nil, fmt.Errorf("parsing %s template: %w", name, err)
		}

		cache[name] = templateSet
	}

	return cache, nil
}

var templateFuncs = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
