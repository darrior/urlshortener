package repository

import (
	"context"
	"sync"

	rmodels "github.com/darrior/urlshortener/internal/repository/models"
)

type MapRepository struct {
	lock sync.Mutex
	urls urlStorage
}

func NewMapRepository() *MapRepository {
	return &MapRepository{
		lock: sync.Mutex{},
		urls: urlStorage{},
	}
}

func (m *MapRepository) AddURL(_ context.Context, userID, id, url string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.urls[id] = record{
		OriginalURL: url,
		UserID:      userID,
	}

	return nil
}

func (m *MapRepository) AddURLs(_ context.Context, userID string, batchURLs rmodels.BatchURLs) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, url := range batchURLs {
		m.urls[url.ID] = record{
			OriginalURL: url.URL,
			UserID:      userID,
		}
	}

	return nil
}

func (m *MapRepository) RemoveURLs(_ context.Context, ids <-chan rmodels.BatchIDsEntry) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for {
		select {
		case id, ok := <-ids:
			if !ok {
				return nil
			}
			if r, ok := m.urls[id.ID]; ok && r.UserID == id.UserID {
				r.Deleted = true
				m.urls[id.ID] = r
			}

		default:
			return nil
		}
	}
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

	if url.Deleted {
		return "", ErrorDeleted
	}

	return url.OriginalURL, nil
}

func (m *MapRepository) GetUserURLs(_ context.Context, userID string) (rmodels.BatchURLs, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var urls rmodels.BatchURLs

	for id, record := range m.urls {
		if record.UserID == userID {
			urls = append(urls, rmodels.BatchURLsEntry{ID: id, URL: record.OriginalURL})
		}
	}

	return urls, nil
}

func (m *MapRepository) Ping(_ context.Context) error {
	return nil
}

func (m *MapRepository) Close() error {
	return nil
}
