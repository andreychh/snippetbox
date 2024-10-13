package storage

import (
	"errors"
	"io"

	"github.com/andreychh/snippetbox/internal/domain"
)

var (
	ErrNoRecord = errors.New("domain: no matching record found")
)

type Storage interface {
	io.Closer
	Snippets() SnippetStorage
}

type SnippetStorage interface {
	Add(snippet *domain.Snippet) error
	GetByID(id int64) (domain.Snippet, error)
	RemoveByID(id int64) error
	Latest() ([]domain.Snippet, error)
}
