package main

import (
	"github.com/andreychh/snippetbox/internal/app"
	cfg "github.com/andreychh/snippetbox/internal/config"
	log "github.com/andreychh/snippetbox/internal/logger"
	srv "github.com/andreychh/snippetbox/internal/server"
	"github.com/andreychh/snippetbox/internal/session"
	"github.com/andreychh/snippetbox/internal/storage/mysql"
	"github.com/andreychh/snippetbox/internal/template"
)

func main() {
	flags := cfg.ParseFlags()
	config, err := cfg.New(flags.YAMLPath)
	if err != nil {
		log.Default().Error("failed to initialize configuration", log.Error(err))
		return
	}

	logger, err := log.New(config.Logger)
	if err != nil {
		log.Default().Error("failed to initialize logger", log.Error(err))
		return
	}

	storage, err := mysql.New(config.Database)
	if err != nil {
		logger.Error("failed to initialize storage", log.Error(err))
		return
	}
	defer func() {
		if err := storage.Close(); err != nil {
			logger.Error("failed to close database connection", log.Error(err))
		}
	}()

	templateRenderer, err := template.NewRenderer()
	if err != nil {
		logger.Error("failed to initialize template renderer", log.Error(err))
		return
	}

	sessionManager := session.NewManager(config.Session, storage)
	application := app.New(logger, storage, templateRenderer, sessionManager)

	server := srv.New(config.Server, application.Routes(), logger.Handler())

	logger.Info("starting server", "address", server.Addr)
	if err = server.ListenAndServeTLS(); err != nil {
		logger.Error("failed to start server", log.Error(err))
		return
	}
}
