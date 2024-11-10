package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/andreychh/snippetbox/internal/template"
)

func (a *App) internalServerError(writer http.ResponseWriter, request *http.Request, err error) {
	a.logger.Error("internal server error",
		slog.String("error", err.Error()),
		slog.String("method", request.Method),
		slog.String("uri", request.RequestURI),
		slog.String("trace", string(debug.Stack())),
	)
	http.Error(
		writer,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (a *App) notFound(writer http.ResponseWriter, request *http.Request) {
	a.logger.Error("page not found",
		slog.String("method", request.Method),
		slog.String("uri", request.RequestURI),
	)

	err := a.writeResponse(writer, template.PageNotFound, template.NewData(), http.StatusNotFound)
	if err != nil {
		err = fmt.Errorf("writing response: %w", err)
		a.internalServerError(writer, request, err)
		return
	}
}
