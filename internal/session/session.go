package session

import (
	"github.com/alexedwards/scs/v2"
	"github.com/andreychh/snippetbox/internal/storage"
	"time"
)

func NewManager(storage storage.Storage) *scs.SessionManager {
	var sessionManager = scs.New()
	sessionManager.Store = storage.Sessions()
	sessionManager.Lifetime = 12 * time.Hour
	return sessionManager
}
