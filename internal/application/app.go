package application

import (
	"log/slog"
	"net/http"

	"github.com/andreychh/snippetbox/internal/storage"
	"github.com/andreychh/snippetbox/internal/templates"

	"github.com/justinas/alice"
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

func (a *App) Routes() http.Handler {
	var mux = http.NewServeMux()
	var fileServer = http.FileServer(http.Dir("./web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", a.notFound)
	mux.HandleFunc("/{$}", a.home)
	mux.HandleFunc("/snippet/view/{id}", a.snippetView)
	mux.HandleFunc("GET /snippet/create", a.snippetCreate)
	mux.HandleFunc("POST /snippet/create", a.snippetCreatePost)

	var middlewareChain = alice.New(a.recoverPanic, a.logRequest, a.addHeaders)

	return middlewareChain.Then(mux)
}
