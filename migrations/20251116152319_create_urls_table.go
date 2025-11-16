// Package migrations provides incremental migrations for DB.
package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUrlsTable, downCreateUrlsTable)
}

func upCreateUrlsTable(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "CREATE TABLE urls (id text NOT NULL, url text, PRIMARY KEY(id))"); err != nil {
		return err
	}

	return nil
}

func downCreateUrlsTable(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS urls"); err != nil {
		return err
	}

	return nil
}
