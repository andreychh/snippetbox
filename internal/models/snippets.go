// #4.5
package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Snippet #4.5 Представляет отдельный snippet
type Snippet struct {
	ID      int64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel #4.5 инкапсулирует пул подключений
type SnippetModel struct {
	DB *sql.DB
}

// Insert #4.6 помещает новый snippet в базу данных
func (m *SnippetModel) Insert(title string, content string, expires int) (int64, error) {
	var stmt = `
INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, fmt.Errorf("executing statement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting id: %w", err)
	}

	return id, nil
}

// Get #4.6 возвращает snippet по его ID
func (m *SnippetModel) Get(id int) (Snippet, error) {
	var stmt = `
SELECT id, title, content, created, expires
FROM snippets
WHERE expires > UTC_TIMESTAMP()
  AND id = ?
`
	var row = m.DB.QueryRow(stmt, id)
	var s Snippet
	var err = row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}
		return Snippet{}, err
	}

	return s, nil
}

// Latest #4.7 возвращает 10 последних добавленных snippet (при условии, что срок их действия не истек)
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `
SELECT id, title, content, created, expires
FROM snippets
WHERE expires > UTC_TIMESTAMP()
ORDER BY id DESC
LIMIT 10
`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// #4.7 *Rows требует закрытия
	// Пока *Rows открыт, он будет поддерживать открытым соединение с базой данных
	defer rows.Close()

	var snippets []Snippet

	// Итерирует по строкам, освобождает ресурсы при завершении
	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// Извлекает любые ошибки, произошедшие во время итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
