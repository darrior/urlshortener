package repository

import (
	"context"
	"database/sql"

	"github.com/darrior/urlshortener/internal/repository/migrations"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) (*DBRepository, error) {
	if err := migrations.Up(context.TODO(), db); err != nil {
		return nil, err
	}

	return &DBRepository{
		db: db,
	}, nil
}

// AddURL implements Repository.
func (d *DBRepository) AddURL(ctx context.Context, id, url string) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO urls (id, url) VALUES ($1, $2)", id, url)
	if err != nil {
		return err
	}

	return nil
}

// GetURL implements Repository.
func (d *DBRepository) GetURL(ctx context.Context, id string) (string, error) {
	row := d.db.QueryRowContext(ctx, "SELECT url FROM urls WHERE id = $1", id)

	var url string
	if err := row.Scan(&url); err != nil {
		return "", err
	}

	return url, nil
}

func (d *DBRepository) Ping(ctx context.Context) error {
	if err := d.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

// Close implements Repository.
func (d *DBRepository) Close() (err error) {
	return d.db.Close()
}
