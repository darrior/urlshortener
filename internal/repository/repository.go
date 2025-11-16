// Package repository provides interface over data sources.
package repository

import (
	"context"
	"errors"
)

type Repository interface {
	AddURL(ctx context.Context, id, url string) (err error)
	GetURL(ctx context.Context, id string) (url string, err error)
	Close() (err error)
}

var ErrorNotFound = errors.New("not found")
