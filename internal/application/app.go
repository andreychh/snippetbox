package application

import (
	"log/slog"
	"net/http"

	"github.com/andreychh/snippetbox/internal/storage"
	"github.com/andreychh/snippetbox/internal/templates"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/alice"
)

type App struct {
	logger           *slog.Logger
	storage          storage.Storage
	templateRenderer *templates.Renderer
	sessionManager   *scs.SessionManager
}

func New(
	logger *slog.Logger,
	storage storage.Storage,
	templateRenderer *templates.Renderer,
	sessionManager *scs.SessionManager,
) *App {
	return &App{
		logger:           logger,
		storage:          storage,
		templateRenderer: templateRenderer,
		sessionManager:   sessionManager,
	}
}

func (a *App) Routes() http.Handler {
	var dynamic = alice.New(a.sessionManager.LoadAndSave)
	var general = alice.New(a.recoverPanic, a.logRequest, a.addHeaders)

	var mux = http.NewServeMux()
	var fileServer = http.FileServer(http.Dir("./web/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", dynamic.ThenFunc(a.notFound))
	mux.Handle("/{$}", dynamic.ThenFunc(a.home))
	mux.Handle("/snippet/view/{id}", dynamic.ThenFunc(a.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(a.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(a.snippetCreatePost))

	return general.Then(mux)
}
