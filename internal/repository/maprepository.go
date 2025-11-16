package repository

import (
	"context"
	"sync"
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

func (r *MapRepository) AddURL(_ context.Context, id, url string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.urls[id] = url

	return nil
}

func (r *MapRepository) GetURL(_ context.Context, id string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	url, ok := r.urls[id]
	if !ok {
		return "", ErrorNotFound
	}

	return url, nil
}

func (r *MapRepository) Close() error {
	return nil
}
