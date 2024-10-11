package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"snippetbox/internal/domain"
	"snippetbox/internal/storage"

	_ "github.com/go-sql-driver/mysql"
)

func openDB(dataSourceName string) (*sql.DB, error) {
	var db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		var closeErr = db.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("pinging database: %w, closing connection: %w", err, closeErr)
		}
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return db, err
}

type Storage struct {
	db       *sql.DB
	snippets SnippetStorage
}

func (s Storage) Snippets() storage.SnippetStorage {
	return s.snippets
}

func New(dataSourceName string) (Storage, error) {
	var db, err = openDB(dataSourceName)
	if err != nil {
		return Storage{}, err
	}

	return Storage{
		db:       db,
		snippets: SnippetStorage{db: db},
	}, nil
}

func (s Storage) Close() error {
	return s.db.Close()
}

type SnippetStorage struct {
	db *sql.DB
}

func (s SnippetStorage) Add(snippet *domain.Snippet) error {
	var query = `
		INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, ?, ?)
	`
	var result, err = s.db.Exec(query, snippet.Title, snippet.Content, snippet.Created, snippet.Expires)
	if err != nil {
		return fmt.Errorf("inserting snippet into database: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("retrieving last insert id: %w", err)
	}

	snippet.ID = id
	return nil
}

func (s SnippetStorage) GetByID(id int64) (domain.Snippet, error) {
	var query = `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP()
		  AND id = ?
	`
	var row = s.db.QueryRow(query, id)
	var snippet, err = s.scanSnippet(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Snippet{}, storage.ErrNoRecord
		}
		return domain.Snippet{}, err
	}

	return snippet, nil
}

func (s SnippetStorage) RemoveByID(id int64) error {
	query := `
		DELETE FROM snippets WHERE id = ?
	`
	var _, err = s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("deleting snippet: %w", err)
	}

	return nil
}

func (s SnippetStorage) Latest() ([]domain.Snippet, error) {
	var query = `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP()
		ORDER BY id DESC
		LIMIT 10
	`
	var rows, err = s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying snippets: %w", err)
	}
	defer rows.Close()

	snippets, err := s.scanSnippets(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning snippets: %w", err)
	}

	return snippets, nil
}

func (s SnippetStorage) scanSnippet(row *sql.Row) (domain.Snippet, error) {
	var snippet domain.Snippet

	var err = row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		return domain.Snippet{}, fmt.Errorf("scanning snippet: %w", err)
	}

	return snippet, nil
}

func (s SnippetStorage) scanSnippets(rows *sql.Rows) ([]domain.Snippet, error) {
	var snippets []domain.Snippet
	var snippet domain.Snippet

	for rows.Next() {
		var err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, fmt.Errorf("scanning snippet: %w", err)
		}
		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
