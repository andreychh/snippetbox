package domain

import (
	"time"
)

type Snippet struct {
	ID      int64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

func NewSnippet(title string, content string, expirationDays int) Snippet {
	var currentTime = time.Now().UTC()
	var expirationTime = currentTime.AddDate(0, 0, expirationDays)

	return Snippet{
		Title:   title,
		Content: content,
		Created: currentTime,
		Expires: expirationTime,
	}
}
