package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/darrior/urlshortener/internal/models"
	"github.com/darrior/urlshortener/internal/repository/migrations"
	"github.com/rs/zerolog/log"
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

func (d *DBRepository) AddURL(ctx context.Context, id, url string) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO urls (id, url) VALUES ($1, $2)", id, url)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBRepository) AddURLs(ctx context.Context, batchURLs models.BatchURLs) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (id, url) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Error().Err(err).Msg("Can not close STMT properly")
		}
	}()

	var errs []error
	for _, url := range batchURLs {
		_, err := stmt.ExecContext(ctx, url.ID, url.URL)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *DBRepository) Count(ctx context.Context) (int, error) {
	row := d.db.QueryRowContext(ctx, "SELECT COUNT (*) FROM urls")
	var count int

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

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

func (d *DBRepository) Close() (err error) {
	return d.db.Close()
}
