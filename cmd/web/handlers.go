package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
)

func (a *App) home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Server", "Go")

	snippets, err := a.storage.Latest()
	if err != nil {
		a.serverError(writer, request, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(writer, "%+v\n", snippet)
	}

	// var files = []string{
	// 	"./ui/html/base.gohtml",
	// 	"./ui/html/partials/nav.gohtml",
	// 	"./ui/html/pages/home.gohtml",
	// }
	//
	// templateSet, err := template.ParseFiles(files...)
	// if err != nil {
	// 	a.serverError(writer, request, fmt.Errorf("parsing template: %w", err))
	// 	return
	// }
	//
	// err = templateSet.ExecuteTemplate(writer, "base", nil)
	// if err != nil {
	// 	a.serverError(writer, request, fmt.Errorf("executing template: %w", err))
	// 	return
	// }
}

func (a *App) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	snippet, err := a.storage.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(writer, request)
		} else {
			a.serverError(writer, request, err)
		}
		return
	}

	fmt.Fprintf(writer, "%+v", snippet)
}

func (a *App) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Display a form for creating a new snippet...\n"))
}

func (a *App) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	// Create some variables holding dummy data. We'll remove these later on during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := a.storage.Insert(title, content, expires)
	if err != nil {
		a.serverError(writer, request, err)
		return
	}

	// #4.6 Переадресация пользователя на страницу с созданным snippet
	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
