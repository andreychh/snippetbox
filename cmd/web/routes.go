package main

import (
	"net/http"
	"snippetbox/internal/config"
)

// #3.5 routes Изолирует создание маршрутизатора
func (a *App) routes(cfg config.Config) *http.ServeMux {
	var mux = http.NewServeMux()
	var fileServer = http.FileServer(http.Dir(cfg.App.StaticDir))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/{$}", a.home)
	mux.HandleFunc("/snippet/view/{id}", a.snippetView)
	mux.HandleFunc("GET /snippet/create", a.snippetCreate)
	mux.HandleFunc("POST /snippet/create", a.snippetCreatePost)

	return mux
}
