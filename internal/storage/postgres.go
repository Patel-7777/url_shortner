package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type URLRecord struct {
	ID          int64
	ShortCode   string
	OriginalURL string
	CreatedAt   time.Time
	VisitCount  int64
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(host, port, user, password, dbname, sslmode string) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreateURL(shortCode, originalURL string) error {
	query := `
		INSERT INTO urls (short_code, original_url, created_at, visit_count)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.db.Exec(query, shortCode, originalURL, time.Now(), 0)
	return err
}

func (s *PostgresStorage) GetURL(shortCode string) (*URLRecord, error) {
	query := `
		SELECT id, short_code, original_url, created_at, visit_count
		FROM urls
		WHERE short_code = $1
	`
	record := &URLRecord{}
	err := s.db.QueryRow(query, shortCode).Scan(
		&record.ID,
		&record.ShortCode,
		&record.OriginalURL,
		&record.CreatedAt,
		&record.VisitCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *PostgresStorage) IncrementVisitCount(shortCode string) error {
	query := `
		UPDATE urls
		SET visit_count = visit_count + 1
		WHERE short_code = $1
	`
	_, err := s.db.Exec(query, shortCode)
	return err
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
