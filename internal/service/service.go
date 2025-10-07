// Package service implemets core logic of URL shortener
package service

import (
	"errors"
	"fmt"

	"github.com/darrior/urlshortener/internal/repository"
)

var (
	ErrorUnknownURL   = errors.New("unkonwn URL")
	ErrorCannotAddURL = errors.New("cannot add URL")
)

type Service struct {
	data repository.Repository
}

func NewService(data repository.Repository) *Service {
	return &Service{
		data: data,
	}
}

func (s *Service) AddURL(url string) (string, error) {
	id := generateURLID()
	for _, err := s.data.GetURL(id); err == nil; {
		id = generateURLID()
	}

	if err := s.data.AddURL(id, url); err != nil {
		return "", fmt.Errorf("%s: %w", ErrorCannotAddURL.Error(), err)
	}

	return id, nil
}

func (s *Service) GetURL(id string) (string, error) {
	url, err := s.data.GetURL(id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrorUnknownURL.Error(), err)
	}

	return url, nil
}
