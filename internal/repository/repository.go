// Package repository provides interface over data sources.
package repository

import (
	"context"
	"errors"

	"github.com/darrior/urlshortener/internal/models"
)

type Repository interface {
	AddURL(ctx context.Context, id, url string) (err error)
	AddURLs(ctx context.Context, batchURLs models.BatchURLs) (err error)
	Count(ctx context.Context) (count int, err error)
	GetURL(ctx context.Context, id string) (url string, err error)
	Ping(ctx context.Context) (err error)
	Close() (err error)
}

var ErrorNotFound = errors.New("not found")
