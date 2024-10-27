package domain

import (
	"fmt"
	"net/http"
	"strconv"
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

func ParseSnippetID(request *http.Request) (int64, error) {
	var idStr = request.PathValue("id")
	var id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing id '%s': %w", idStr, err)
	}
	if id < 1 {
		return 0, fmt.Errorf("snippet id '%d' must be a positive integer", id)
	}
	return id, nil
}
