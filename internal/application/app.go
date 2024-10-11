package application

import (
	"log/slog"
	"net/http"
	"snippetbox/internal/storage"
	"snippetbox/internal/templates"
)

type App struct {
	logger           *slog.Logger
	storage          storage.Storage
	templateRenderer templates.Renderer
}

func New(logger *slog.Logger, storage storage.Storage, templateRenderer templates.Renderer) App {
	return App{
		logger:           logger,
		storage:          storage,
		templateRenderer: templateRenderer,
	}
}

func (a *App) Routes() *http.ServeMux {
	var mux = http.NewServeMux()
	var fileServer = http.FileServer(http.Dir("./web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", a.home)
	mux.HandleFunc("/snippet/view/{id}", a.snippetView)
	mux.HandleFunc("GET /snippet/create", a.snippetCreate)
	mux.HandleFunc("POST /snippet/create", a.snippetCreatePost)

	return mux
}
