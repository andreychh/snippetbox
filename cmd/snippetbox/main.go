package main

import (
	"net/http"
	"os"

	"github.com/andreychh/snippetbox/internal/application"
	cfg "github.com/andreychh/snippetbox/internal/config"
	log "github.com/andreychh/snippetbox/internal/logger"
	"github.com/andreychh/snippetbox/internal/session"
	"github.com/andreychh/snippetbox/internal/storage/mysql"
	"github.com/andreychh/snippetbox/internal/template"
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

	templateRenderer, err := template.NewRenderer()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var sessionManager = session.NewManager(storage)

	var app = application.New(logger, storage, templateRenderer, sessionManager)

	logger.Info("starting server", "address", config.App.Addr())
	err = http.ListenAndServe(config.App.Addr(), app.Routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
