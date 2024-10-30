package application

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/andreychh/snippetbox/internal/domain"
	"github.com/andreychh/snippetbox/internal/storage"
	"github.com/andreychh/snippetbox/internal/template"
)

func (a *App) home(writer http.ResponseWriter, request *http.Request) {
	var snippets, err = a.storage.Snippets().Latest()
	if err != nil {
		err = fmt.Errorf("fetching latest snippets: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	var templateData = template.NewData(template.WithSnippets(snippets))
	pageContent, err := a.templateRenderer.RenderPage(template.PageHome, templateData)
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
	var id, err = domain.ParseSnippetID(request)
	if err != nil {
		err = fmt.Errorf("parsing snippet ID: %w", err)
		a.logger.Error("error occurred", slog.String("error", err.Error()))
		a.notFound(writer, request)
		return
	}

	snippet, err := a.storage.Snippets().GetByID(id)
	if err != nil {
		err = fmt.Errorf("fetching snippet by ID %d: %w", id, err)
		if errors.Is(err, storage.ErrNoRecord) {
			a.logger.Error("error occurred", slog.String("error", err.Error()))
			a.notFound(writer, request)
		} else {
			a.internalServerError(writer, request, err)
		}
		return
	}

	var flash = a.sessionManager.PopString(request.Context(), "flash")
	var templateData = template.NewData(template.WithSnippet(snippet), template.WithFlash(flash))
	pageContent, err := a.templateRenderer.RenderPage(template.PageView, templateData)
	if err != nil {
		err = fmt.Errorf("rendering snippet-view page: %w", err)
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
	var templateData = template.NewData(template.WithForm(domain.Form{Expires: 365}))
	var pageContent, err = a.templateRenderer.RenderPage(template.PageCreate, templateData)
	if err != nil {
		err = fmt.Errorf("rendering snippet-create page: %w", err)
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

func (a *App) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	var form, err = domain.ParseSnippetCreateForm(request)
	if err != nil {
		err = fmt.Errorf("parsing form: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	if !form.Valid() {
		var templateData = template.NewData(template.WithForm(form))
		var pageContent, err = a.templateRenderer.RenderPage(template.PageCreate, templateData)
		if err != nil {
			err = fmt.Errorf("rendering snippet-create page: %w", err)
			a.internalServerError(writer, request, err)
			return
		}

		writer.WriteHeader(http.StatusUnprocessableEntity)
		n, err := writer.Write(pageContent)
		if err != nil {
			err = fmt.Errorf("writing response (bytes written: %d): %w", n, err)
			a.internalServerError(writer, request, err)
			return
		}
		return
	}

	var snippet = domain.NewSnippet(form.Title, form.Content, form.Expires)
	err = a.storage.Snippets().Add(&snippet)
	if err != nil {
		err = fmt.Errorf("adding snippet: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	a.sessionManager.Put(request.Context(), "flash", "Snippet successfully created!")
	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}
