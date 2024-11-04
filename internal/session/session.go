package session

import (
	cfg "github.com/andreychh/snippetbox/internal/config"
	"github.com/andreychh/snippetbox/internal/storage"

	"github.com/alexedwards/scs/v2"
)

func NewManager(config cfg.Session, storage storage.Storage) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = storage.Sessions()
	sessionManager.Lifetime = config.Lifetime
	sessionManager.Cookie.Secure = true
	return sessionManager
}
