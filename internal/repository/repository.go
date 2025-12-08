// Package repository provides interface over data sources.
package repository

import (
	"context"
	"errors"

	rmodels "github.com/darrior/urlshortener/internal/repository/models"
)

type Repository interface {
	AddURL(ctx context.Context, userID, id, url string) (err error)
	AddURLs(ctx context.Context, userID string, batchURLs rmodels.BatchURLs) (err error)
	Count(ctx context.Context) (count int, err error)
	GetURL(ctx context.Context, id string) (url string, err error)
	GetUserURLs(ctx context.Context, userID string) (urls rmodels.BatchURLs, err error)
	Ping(ctx context.Context) (err error)
	Close() (err error)
}

var ErrorNotFound = errors.New("not found")
