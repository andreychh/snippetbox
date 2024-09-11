package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (a *App) home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Server", "Go")

	var files = []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	var templateSet, err = template.ParseFiles(files...)
	if err != nil {
		a.serverError(writer, request, fmt.Errorf("parsing template: %w", err)) // #3.4
		return
	}

	err = templateSet.ExecuteTemplate(writer, "base", nil)
	if err != nil {
		a.serverError(writer, request, fmt.Errorf("executing template: %w", err)) // #3.4
		return
	}
}

func (a *App) snippetView(writer http.ResponseWriter, request *http.Request) {
	var snippetId, err = strconv.Atoi(request.PathValue("id"))
	if err != nil || snippetId < 1 {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d...\n", snippetId)
}

func (a *App) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Display a form for creating a new snippet...\n"))
}

func (a *App) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte("Save a new snippet...\n"))
}
