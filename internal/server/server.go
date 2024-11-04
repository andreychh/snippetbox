package server

import (
	"crypto/tls"
	"log/slog"
	"net/http"

	cfg "github.com/andreychh/snippetbox/internal/config"
)

type Server struct {
	*http.Server
	certFilepath string
	keyFilepath  string
}

func New(config cfg.Server, httpHandler http.Handler, slogHandler slog.Handler) *Server {
	errorLog := slog.NewLogLogger(slogHandler, config.LogLevel)
	tlsConfig := &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}

	server := &http.Server{
		Addr:         config.Address(),
		Handler:      httpHandler,
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return &Server{server, config.CertPath, config.KeyPath}
}

func (s *Server) ListenAndServeTLS() error {
	return s.Server.ListenAndServeTLS(s.certFilepath, s.keyFilepath)
}
