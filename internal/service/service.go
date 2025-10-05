// Package service implemets core logic of URL shortener
package service

import "errors"

type Service struct {
	urls map[string]string
}

func NewService() *Service {
	return &Service{
		urls: make(map[string]string),
	}
}

func (s *Service) AddURL(url string) string {
	shortURL := generateURL()
	for _, ok := s.urls[shortURL]; ok; {
		shortURL = generateURL()
	}

	s.urls[shortURL] = url

	return shortURL
}

func (s *Service) GetURL(shortURL string) (string, error) {
	url, ok := s.urls[shortURL]

	if !ok {
		return "", errors.New("URL not found")
	}
	return url, nil
}
