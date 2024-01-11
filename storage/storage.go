package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dbStore struct {
	pool *pgxpool.Pool
}

func (s *dbStore) GetURL(ctx context.Context, shortURL string) (string, error) {
	rows, err := s.pool.Query(ctx, `SELECT url FROM links WHERE short_url = $1`, shortURL);
	if err != nil {
		return "", err
	}

	var url string

	for rows.Next() {
		err := rows.Scan(&url);

		if err != nil {
			return "", err
		}
	}

	if url == "" {
		return "", fmt.Errorf("url with %s Short URL doesn't exist", shortURL)
	}

	return url, nil
} 

func (s *dbStore) GetShortURL(ctx context.Context, url string) (string, error) {
	rows, err := s.pool.Query(ctx, `SELECT short_url FROM links WHERE url = $1`, url);
	if err != nil {
		return "", err
	}

	var shortURL string

	for rows.Next() {
		err := rows.Scan(&shortURL);

		if err != nil {
			return "", err
		}
	}

	if shortURL == "" {
		return "", fmt.Errorf("url with %s Short URL doesn't exist", url)
	}

	return shortURL, nil
} 

func (s *dbStore) AddURL(ctx context.Context, url, shortURL string) error {

	_, err := s.pool.Exec(ctx, `INSERT INTO links(url, short_url) values ($1, $2)`, url, shortURL)
	if err != nil {
		return err
	}

	return nil
}

func (s *dbStore) HaveURL(ctx context.Context, url string) (bool, error) {

	rows, err := s.pool.Query(ctx, `SELECT EXISTS(SELECT * FROM links WHERE url=$1)`, url)
	if err != nil {
		return false, err
	}

	var exists bool

	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			return false, err
		}
	}

	return exists, nil
}

func (s *dbStore) CreateSchema() error {
	query := `CREATE TABLE IF NOT EXISTS links (
		short_url	text PRIMARY KEY,
		url		    text
	);`
	_, err := s.pool.Exec(context.Background(), query);

	if err != nil {
		return err
	}
	return nil
}

func NewPostgresStore(ctx context.Context) (*dbStore, error) {
	db_user := os.Getenv("db_user")
	db_pass := os.Getenv("db_pass")
	db_name := os.Getenv("db_name")

	connStr := "postgres://"+db_user+":"+db_pass+"@localhost/"+db_name+"?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr);
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &dbStore{pool: pool}, nil
}