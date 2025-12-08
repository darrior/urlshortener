package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/darrior/urlshortener/internal/repository/migrations"
	rmodels "github.com/darrior/urlshortener/internal/repository/models"
	"github.com/rs/zerolog/log"
)

type ErrorURLExists struct {
	ID  string
	URL string
}

func newErrorIDExists(id, url string) error {
	return &ErrorURLExists{
		ID:  id,
		URL: url,
	}
}

func (e *ErrorURLExists) Error() string {
	return fmt.Sprintf("url %s exists with id %s", e.URL, e.ID)
}

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

func (d *DBRepository) AddURL(ctx context.Context, userID, id, url string) error {
	row := d.db.QueryRowContext(ctx, "INSERT INTO urls (id, url, users) VALUES ($1, $2, $3) ON CONFLICT (url) DO UPDATE SET users = urls.users || EXCLUDED.users RETURNING id", id, url, []string{userID})
	var inserted string
	if err := row.Scan(&inserted); err != nil {
		log.Error().Err(err).Msg("Can not scan row")
		return fmt.Errorf("can not parse row: %w", err)
	}

	log.Info().Str("inserted", inserted).Msg("DB returns value")
	if inserted != id {
		return newErrorIDExists(inserted, url)
	}
	return nil
}

func (d *DBRepository) AddURLs(ctx context.Context, userID string, batchURLs rmodels.BatchURLs) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("can not begin transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (id, url, users) VALUES ($1, $2, $3) ON CONFLICT (url) DO UPDATE SET users = urls.users || EXCLUDED.users RETURNING id")
	if err != nil {
		return fmt.Errorf("can not prepare query: %w", err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Error().Err(err).Msg("Can not close STMT properly")
		}
	}()

	var errs []error
	for _, url := range batchURLs {
		row := stmt.QueryRowContext(ctx, url.ID, url.URL, []string{userID})

		var inserted string
		if err := row.Scan(&inserted); err != nil {
			return fmt.Errorf("can not parse row: %w", err)
		}

		if url.ID != inserted {
			errs = append(errs, newErrorIDExists(url.ID, url.URL))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can not commit transaction: %w", err)
	}

	return nil
}

func (d *DBRepository) Count(ctx context.Context) (int, error) {
	row := d.db.QueryRowContext(ctx, "SELECT COUNT (*) FROM urls")
	var count int

	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("can not parse row: %w", err)
	}

	return count, nil
}

func (d *DBRepository) GetURL(ctx context.Context, id string) (string, error) {
	row := d.db.QueryRowContext(ctx, "SELECT url FROM urls WHERE id = $1", id)

	var url string
	if err := row.Scan(&url); err != nil {
		return "", fmt.Errorf("can not parse row: %w", err)
	}

	return url, nil
}

func (d *DBRepository) GetUserURLs(ctx context.Context, userID string) (rmodels.BatchURLs, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, url FROM urls WHERE $1 = ANY(users)", userID)
	if err != nil {
		return rmodels.BatchURLs{}, fmt.Errorf("can not get user's urls: %w", err)
	}

	var (
		errs []error
		urls rmodels.BatchURLs
	)

	for rows.Next() {
		var id, url string
		if err := rows.Scan(&id, &url); err != nil {
			errs = append(errs, err)
			continue
		}

		urls = append(urls, rmodels.BatchURLEntry{
			ID:  id,
			URL: url,
		})

	}
	if len(errs) != 0 {
		return rmodels.BatchURLs{}, errors.Join(errs...)
	}

	if err := rows.Err(); err != nil {
		return rmodels.BatchURLs{}, fmt.Errorf("rows contains error: %w", err)
	}

	return urls, nil
}

func (d *DBRepository) Ping(ctx context.Context) error {
	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("can not ping DB: %w", err)
	}

	return nil
}

func (d *DBRepository) Close() (err error) {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("can not close db: %w", err)
	}

	return nil
}
