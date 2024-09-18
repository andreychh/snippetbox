package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// 3.4 serverError логирует ошибку и возвращает пользователю 500 Internal Server Error
func (a *App) serverError(writer http.ResponseWriter, request *http.Request, err error) {
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

func (a *App) clientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}
