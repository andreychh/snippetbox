package application

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (a *App) internalServerError(writer http.ResponseWriter, request *http.Request, err error) {
	a.logger.Error(
		err.Error(),
		slog.String("method", request.Method),
		slog.String("uri", request.URL.RequestURI()),
		slog.String("trace", string(debug.Stack())),
	)
	http.Error(
		writer,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func (a *App) notFound(writer http.ResponseWriter, request *http.Request, err error) {
	a.logger.Error(
		err.Error(),
		slog.String("method", request.Method),
		slog.String("uri", request.URL.RequestURI()),
		slog.String("trace", string(debug.Stack())),
	)
	http.Error(
		writer,
		http.StatusText(http.StatusNotFound),
		http.StatusNotFound,
	)
}
