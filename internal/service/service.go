// Package service implemets core logic of URL shortener
package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/darrior/urlshortener/internal/models/api"
	"github.com/darrior/urlshortener/internal/repository"
	rmodels "github.com/darrior/urlshortener/internal/repository/models"
	"github.com/darrior/urlshortener/internal/service/auth"
	"github.com/rs/zerolog/log"
)

var (
	ErrorUnknownURL   = errors.New("unkonwn URL")
	ErrorCannotAddURL = errors.New("cannot add URL")
	ErrorURLExists    = errors.New("passed existing URL")
)

type IService interface {
	auth.Auth
	AddURL(ctx context.Context, userID, longURL string) (shortURL string, err error)
	AddURLs(ctx context.Context, userID string, longURLs api.ShortenerBatchRequest) (shortURLs api.ShortenerBatchResponse, err error)
	GetURL(ctx context.Context, id string) (longURL string, err error)
	GetUserURLs(ctx context.Context, userID string) (urls api.UserURLsResponse, err error)
	Ping(ctx context.Context) (err error)
}

type Service struct {
	auth.Auth
	lock        sync.RWMutex
	data        repository.Repository
	baseAddress string
}

func NewService(data repository.Repository, baseAddress string, authKey string) *Service {
	return &Service{
		Auth:        auth.NewHS256Auth(authKey),
		data:        data,
		baseAddress: baseAddress,
	}
}

func (s *Service) AddURL(ctx context.Context, userID, longURL string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	count, err := s.data.Count(ctx)
	if err != nil {
		return "", err
	}

	id := generateURLID(count)

	if err := s.data.AddURL(ctx, userID, id, longURL); err != nil {
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

func (s *Service) AddURLs(ctx context.Context, userID string, longURLs api.ShortenerBatchRequest) (api.ShortenerBatchResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	count, err := s.data.Count(ctx)
	if err != nil {
		return api.ShortenerBatchResponse{}, err
	}

	var res api.ShortenerBatchResponse
	var urls rmodels.BatchURLs
	for _, entry := range longURLs {
		id := generateURLID(count)
		shortURL, err := url.JoinPath(s.baseAddress, id)
		if err != nil {
			return api.ShortenerBatchResponse{}, err
		}

		res = append(res, api.BatchResponseEntry{
			CorrelationID: entry.CorrelationID,
			ShortURL:      shortURL,
		})
		urls = append(urls, rmodels.BatchURLEntry{
			ID:  id,
			URL: entry.OriginalURL,
		})

		count += 1
	}

	if err := s.data.AddURLs(ctx, userID, urls); err != nil {
		return api.ShortenerBatchResponse{}, err
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

func (s *Service) GetUserURLs(ctx context.Context, userID string) (api.UserURLsResponse, error) {
	urls, err := s.data.GetUserURLs(ctx, userID)
	if err != nil {
		return api.UserURLsResponse{}, fmt.Errorf("can not get user URLs from repository: %w", err)
	}

	var (
		userURLs api.UserURLsResponse
		errs     []error
	)
	for _, batchURL := range urls {
		shortURL, err := url.JoinPath(s.baseAddress, batchURL.ID)
		if err != nil {
			errs = append(errs, err)
		}

		userURLs = append(userURLs, api.UserURLsResponseEntry{
			ShortURL:    shortURL,
			OriginalURL: batchURL.URL,
		})
	}

	if len(errs) != 0 {
		return api.UserURLsResponse{}, fmt.Errorf("can not get urls: %w", errors.Join(errs...))
	}

	return userURLs, nil
}

func (s *Service) Ping(ctx context.Context) error {
	if err := s.data.Ping(ctx); err != nil {
		return err
	}

	return nil
}
