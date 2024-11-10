package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/andreychh/snippetbox/internal/domain"
	log "github.com/andreychh/snippetbox/internal/logger"
	"github.com/andreychh/snippetbox/internal/storage"
	"github.com/andreychh/snippetbox/internal/template"
)

func (a *App) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := a.storage.Snippets().Latest()
	if err != nil {
		err = fmt.Errorf("fetching latest snippets: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	data := template.NewData(template.WithSnippets(snippets))
	err = a.writeResponse(writer, template.PageHome, data, http.StatusOK)
	if err != nil {
		err = fmt.Errorf("writing response: %w", err)
		a.internalServerError(writer, request, err)
		return
	}
}

func (a *App) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := domain.ParseSnippetID(request)
	if err != nil {
		a.logger.Error("error occurred", log.Error(fmt.Errorf("parsing snippet ID: %w", err)))
		a.notFound(writer, request)
		return
	}

	snippet, err := a.storage.Snippets().GetByID(id)
	if err != nil {
		err = fmt.Errorf("fetching snippet by ID %d: %w", id, err)
		if errors.Is(err, storage.ErrNoRecord) {
			a.logger.Error("error occurred", log.Error(err))
			a.notFound(writer, request)
		} else {
			a.internalServerError(writer, request, err)
		}
		return
	}

	flash := a.sessionManager.PopString(request.Context(), "flash")
	data := template.NewData(template.WithSnippet(snippet), template.WithFlash(flash))
	err = a.writeResponse(writer, template.PageView, data, http.StatusOK)
	if err != nil {
		err = fmt.Errorf("writing response: %w", err)
		a.internalServerError(writer, request, err)
		return
	}
}

func (a *App) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	data := template.NewData(template.WithForm(domain.Form{Expires: 365}))
	err := a.writeResponse(writer, template.PageCreate, data, http.StatusOK)
	if err != nil {
		err = fmt.Errorf("writing response: %w", err)
		a.internalServerError(writer, request, err)
		return
	}
}

func (a *App) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	form, err := domain.ParseSnippetCreateForm(request)
	if err != nil {
		err = fmt.Errorf("parsing form: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	if !form.Valid() {
		data := template.NewData(template.WithForm(form))
		err = a.writeResponse(writer, template.PageCreate, data, http.StatusUnprocessableEntity)
		if err != nil {
			err = fmt.Errorf("writing response: %w", err)
			a.internalServerError(writer, request, err)
			return
		}
		return
	}

	snippet := domain.NewSnippet(form.Title, form.Content, form.Expires)
	err = a.storage.Snippets().Add(&snippet)
	if err != nil {
		err = fmt.Errorf("adding snippet: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	a.sessionManager.Put(request.Context(), "flash", "Snippet successfully created!")
	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}
