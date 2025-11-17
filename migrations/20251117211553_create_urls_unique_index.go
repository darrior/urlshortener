package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUrlsUniqueIndex, downCreateUrlsUniqueIndex)
}

func upCreateUrlsUniqueIndex(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "CREATE UNIQUE INDEX unique_urls ON urls (url)"); err != nil {
		return err
	}

	return nil
}

func downCreateUrlsUniqueIndex(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "DROP INDEX IF EXISTS unique_urls"); err != nil {
		return err
	}

	return nil
}
