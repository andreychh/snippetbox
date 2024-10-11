package main

import (
	"net/http"
	"os"

	"snippetbox/internal/application"
	cfg "snippetbox/internal/config"
	log "snippetbox/internal/logger"
	"snippetbox/internal/storage/mysql"
	"snippetbox/internal/templates"
)

func main() {
	var flags = cfg.ParseFlags()
	var config = cfg.MustLoad(flags.YAMLPath)

	var logger = log.New(config.EnvName)
	var storage, err = mysql.New(config.Database.DSN())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer storage.Close()

	templateRenderer, err := templates.NewRenderer()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var app = application.New(logger, storage, templateRenderer)

	logger.Info("starting server", "address", config.App.Addr())
	err = http.ListenAndServe(config.App.Addr(), app.Routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
