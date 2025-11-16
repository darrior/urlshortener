// Package migrations provide functions to migrate DB.
package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	_ "github.com/darrior/urlshortener/migrations"
)

func Up(ctx context.Context, db *sql.DB) error {
	err := goose.UpContext(ctx, db, ".")
	if err != nil {
		return err
	}

	return nil
}

func Down(ctx context.Context, db *sql.DB) error {
	err := goose.DownContext(ctx, db, ".")
	if err != nil {
		return err
	}

	return nil
}
