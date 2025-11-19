// Package migrations provide functions to migrate DB.
package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/darrior/urlshortener/migrations"
	"github.com/pressly/goose/v3"
)

func Up(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(migrations.EmbedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("can not set dialect: %w", err)
	}

	err := goose.UpContext(ctx, db, ".")
	if err != nil {
		return err
	}

	return nil
}

func Down(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(migrations.EmbedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("can not set dialect: %w", err)
	}

	err := goose.DownContext(ctx, db, ".")
	if err != nil {
		return err
	}

	return nil
}
