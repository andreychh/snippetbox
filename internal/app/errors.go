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

	var pageContent, err = a.templateRenderer.RenderPage(template.PageNotFound, template.NewData())
	if err != nil {
		err = fmt.Errorf("rendering not-found page: %w", err)
		a.internalServerError(writer, request, err)
		return
	}

	writer.WriteHeader(http.StatusNotFound)
	n, err := writer.Write(pageContent)
	if err != nil {
		err = fmt.Errorf("writing response (bytes written: %d): %w", n, err)
		a.internalServerError(writer, request, err)
		return
	}
}
