package repository

import (
	"context"
	"sync"

	"github.com/darrior/urlshortener/internal/models"
)

type MapRepository struct {
	lock sync.Mutex
	urls map[string]string
}

func NewMapRepository() *MapRepository {
	return &MapRepository{
		lock: sync.Mutex{},
		urls: map[string]string{},
	}
}

func (m *MapRepository) AddURL(_ context.Context, id, url string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.urls[id] = url

	return nil
}

func (m *MapRepository) AddURLs(_ context.Context, batchURLs models.BatchURLs) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, url := range batchURLs {
		m.urls[url.ID] = url.URL
	}

	return nil
}

func (m *MapRepository) Count(_ context.Context) (int, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	return len(m.urls), nil
}

func (m *MapRepository) GetURL(_ context.Context, id string) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	url, ok := m.urls[id]
	if !ok {
		return "", ErrorNotFound
	}

	return url, nil
}

func (m *MapRepository) Ping(_ context.Context) error {
	return nil
}

func (m *MapRepository) Close() error {
	return nil
}
