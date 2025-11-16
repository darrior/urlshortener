// Package service implemets core logic of URL shortener
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/darrior/urlshortener/internal/repository"
)

var (
	ErrorUnknownURL   = errors.New("unkonwn URL")
	ErrorCannotAddURL = errors.New("cannot add URL")
)

type IService interface {
	AddURL(ctx context.Context, longURL string) (shortURL string, err error)
	GetURL(ctx context.Context, id string) (longURL string, err error)
	Ping(ctx context.Context) (err error)
}

type Service struct {
	data        repository.Repository
	baseAddress string
	db          *sql.DB
}

func NewService(data repository.Repository, baseAddress string, db *sql.DB) *Service {
	return &Service{
		data:        data,
		baseAddress: baseAddress,
		db:          db,
	}
}

func (s *Service) AddURL(ctx context.Context, longURL string) (string, error) {
	id := generateURLID()
	for _, err := s.data.GetURL(ctx, id); err == nil; {
		id = generateURLID()
	}

	if err := s.data.AddURL(ctx, id, longURL); err != nil {
		return "", fmt.Errorf("%s: %w", ErrorCannotAddURL.Error(), err)
	}

	shortURL, err := url.JoinPath(s.baseAddress, id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrorCannotAddURL.Error(), err)
	}

	return shortURL, nil
}

func (s *Service) GetURL(ctx context.Context, id string) (string, error) {
	longURL, err := s.data.GetURL(ctx, id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrorUnknownURL.Error(), err)
	}

	return longURL, nil
}

func (s *Service) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
