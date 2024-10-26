package application

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreychh/snippetbox/internal/domain"
	"github.com/andreychh/snippetbox/internal/storage"
	"github.com/andreychh/snippetbox/internal/templates"
)

func (a *App) home(writer http.ResponseWriter, request *http.Request) {
	var snippets, err = a.storage.Snippets().Latest()
	if err != nil {
		err = fmt.Errorf("fetching latest snippets: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	var templateData = templates.NewTemplateData(templates.WithSnippets(snippets))
	pageContent, err := a.templateRenderer.HomePage(templateData)
	if err != nil {
		err = fmt.Errorf("rendering home page: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	n, err := writer.Write(pageContent)
	if err != nil {
		err = fmt.Errorf("writing response (bytes written: %d): %w", n, err)
		a.internalServerError(writer, request, err)
		return
	}
}

func (a *App) snippetView(writer http.ResponseWriter, request *http.Request) {
	var idStr = request.PathValue("id")
	var id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		err = fmt.Errorf("invalid ID: %s. Must be a positive integer", idStr)
		a.notFound(writer, request, err)
		return
	}

	snippet, err := a.storage.Snippets().GetByID(id)
	if err != nil {
		err = fmt.Errorf("fetching snippet by ID %d: %w", id, err)
		if errors.Is(err, storage.ErrNoRecord) {
			a.notFound(writer, request, err)
		} else {
			a.internalServerError(writer, request, err)
		}
		return
	}

	var templateData = templates.NewTemplateData(templates.WithSnippet(snippet))
	pageContent, err := a.templateRenderer.SnippetViewPage(templateData)
	if err != nil {
		err = fmt.Errorf("rendering snippet view page: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	n, err := writer.Write(pageContent)
	if err != nil {
		err = fmt.Errorf("writing response (bytes written: %d): %w", n, err)
		a.internalServerError(writer, request, err)
		return
	}
}

func (a *App) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Display a form for creating a new snippet...\n"))
}

func (a *App) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	var snippet = domain.NewSnippet(title, content, expires)
	var err = a.storage.Snippets().Add(&snippet)
	if err != nil {
		err = fmt.Errorf("adding snippet: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}
