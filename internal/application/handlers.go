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
	writer.Header().Add("Server", "Go")

	var snippets, err = a.storage.Snippets().Latest()
	if err != nil {
		a.internalServerError(writer, request, err)
		return
	}

	var data = templates.NewTemplateData(templates.WithSnippets(snippets))
	err = a.templateRenderer.HomePage(writer, data)
	if err != nil {
		a.internalServerError(writer, request, err)
		return
	}
}

func (a *App) snippetView(writer http.ResponseWriter, request *http.Request) {
	var idStr = request.PathValue("id")
	var id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		a.notFound(writer, request, fmt.Errorf("invalid ID: %s. Must be a positive integer", idStr))
		return
	}

	snippet, err := a.storage.Snippets().GetByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			a.notFound(writer, request, err)
		} else {
			a.internalServerError(writer, request, err)
		}
		return
	}

	var data = templates.NewTemplateData(templates.WithSnippet(snippet))
	err = a.templateRenderer.SnippetViewPage(writer, data)
	if err != nil {
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
		a.internalServerError(writer, request, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}
