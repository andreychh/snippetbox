package domain

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type SnippetCreateForm struct {
	Title   string
	Content string
	Expires int

	FieldErrors map[string]string
}

func NewSnippetCreateForm(title string, content string, expires int) SnippetCreateForm {
	var form = SnippetCreateForm{Title: title, Content: content, Expires: expires}
	form.setFieldErrors()
	return form
}

func (f *SnippetCreateForm) Valid() bool {
	return len(f.FieldErrors) == 0
}

func (f *SnippetCreateForm) setFieldErrors() {
	f.FieldErrors = make(map[string]string)

	if strings.TrimSpace(f.Title) == "" {
		f.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(f.Title) > 100 {
		f.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(f.Content) == "" {
		f.FieldErrors["content"] = "This field cannot be blank"
	}

	if f.Expires != 1 && f.Expires != 7 && f.Expires != 365 {
		f.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}
}

func ParseSnippetCreateForm(request *http.Request) (SnippetCreateForm, error) {
	var err = request.ParseForm()
	if err != nil {
		return SnippetCreateForm{}, fmt.Errorf("parsing form data: %w", err)
	}

	var (
		title      = strings.TrimSpace(request.PostForm.Get("title"))
		content    = strings.TrimSpace(request.PostForm.Get("content"))
		expiresStr = request.PostForm.Get("expires")
	)

	expires, err := strconv.Atoi(expiresStr)
	if err != nil {
		return SnippetCreateForm{}, fmt.Errorf("parsing expiration value '%s': %w", expiresStr, err)
	}

	return NewSnippetCreateForm(title, content, expires), nil
}
