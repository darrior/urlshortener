// Package service implemets core logic of URL shortener
package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/darrior/urlshortener/internal/models"
	"github.com/darrior/urlshortener/internal/repository"
	"github.com/rs/zerolog/log"
)

var (
	ErrorUnknownURL   = errors.New("unkonwn URL")
	ErrorCannotAddURL = errors.New("cannot add URL")
	ErrorURLExists    = errors.New("passed existing URL")
)

type IService interface {
	AddURL(ctx context.Context, longURL string) (shortURL string, err error)
	AddURLs(ctx context.Context, longURLs models.ShortenerBatchRequest) (shortURLs models.ShortenerBatchResponse, err error)
	GetURL(ctx context.Context, id string) (longURL string, err error)
	Ping(ctx context.Context) (err error)
}

type Service struct {
	lock        sync.RWMutex
	data        repository.Repository
	baseAddress string
}

func NewService(data repository.Repository, baseAddress string) *Service {
	return &Service{
		data:        data,
		baseAddress: baseAddress,
	}
}

func (s *Service) AddURL(ctx context.Context, longURL string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	count, err := s.data.Count(ctx)
	if err != nil {
		return "", err
	}

	id := generateURLID(count)

	if err := s.data.AddURL(ctx, id, longURL); err != nil {
		var ue *repository.ErrorURLExists
		if !errors.As(err, &ue) {
			log.Error().Err(err).Msg("error")
			return "", fmt.Errorf("%s: %w", ErrorCannotAddURL.Error(), err)
		}

		log.Info().Str("id", id).Msg("Repository returns existed ID")

		shortURL, err := url.JoinPath(s.baseAddress, ue.ID)
		if err != nil {
			return "", fmt.Errorf("%s: %w", ErrorCannotAddURL.Error(), err)
		}

		return shortURL, ErrorURLExists
	}

	shortURL, err := url.JoinPath(s.baseAddress, id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrorCannotAddURL.Error(), err)
	}

	return shortURL, nil
}

func (s *Service) AddURLs(ctx context.Context, longURLs models.ShortenerBatchRequest) (models.ShortenerBatchResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	count, err := s.data.Count(ctx)
	if err != nil {
		return models.ShortenerBatchResponse{}, err
	}

	var res models.ShortenerBatchResponse
	var urls models.BatchURLs
	for _, entry := range longURLs {
		id := generateURLID(count)
		shortURL, err := url.JoinPath(s.baseAddress, id)
		if err != nil {
			return models.ShortenerBatchResponse{}, err
		}

		res = append(res, models.BatchResponseEntry{
			CorrelationID: entry.CorrelationID,
			ShortURL:      shortURL,
		})
		urls = append(urls, models.BatchURLEntry{
			ID:  id,
			URL: entry.OriginalURL,
		})

		count += 1
	}

	if err := s.data.AddURLs(ctx, urls); err != nil {
		return models.ShortenerBatchResponse{}, err
	}

	return res, nil
}

func (s *Service) GetURL(ctx context.Context, id string) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	longURL, err := s.data.GetURL(ctx, id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrorUnknownURL.Error(), err)
	}

	return longURL, nil
}

func (s *Service) Ping(ctx context.Context) error {
	if err := s.data.Ping(ctx); err != nil {
		return err
	}

	return nil
}
