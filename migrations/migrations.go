// Package migrations provides all embedded SQL migrations in one variable.
package migrations

import "embed"

//go:embed *.sql
var EmbedMigrations embed.FS
