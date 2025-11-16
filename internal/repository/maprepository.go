package repository

import "sync"

type MapRepository struct {
	lock sync.Mutex
	urls map[string]string
}

var _ Repository = (*MapRepository)(nil)

func NewMapRepository() *MapRepository {
	return &MapRepository{
		lock: sync.Mutex{},
		urls: map[string]string{},
	}
}

func (r *MapRepository) AddURL(id, url string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.urls[id] = url

	return nil
}

func (r *MapRepository) GetURL(id string) (string, error) {
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
