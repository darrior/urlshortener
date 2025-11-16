package repository

import (
	"context"
	"os"
	"sync"

	"github.com/darrior/urlshortener/internal/repository/storage"
	"github.com/rs/zerolog/log"
)

type urlStorage map[string]string

type FSRepository struct {
	lock sync.Mutex
	urls urlStorage
	file *os.File
}

func NewFSRepository(file *os.File) (*FSRepository, error) {
	var urls urlStorage
	if err := storage.ReadFile(file, &urls); err != nil {
		log.Warn().Err(err).Msg("Can not read urls from storage file")
		urls = urlStorage{}
	}

	r := &FSRepository{
		lock: sync.Mutex{},
		urls: urls,
		file: file,
	}

	return r, nil
}

func (f *FSRepository) AddURL(_ context.Context, id, url string) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.urls[id] = url
	if err := storage.UpdateFile(f.file, f.urls); err != nil {
		return err
	}

	return nil
}

func (f *FSRepository) GetURL(_ context.Context, id string) (string, error) {
	url, ok := f.urls[id]
	if !ok {
		return "", ErrorNotFound
	}

	return url, nil
}

func (f *FSRepository) Ping(_ context.Context) error {
	return nil
}

func (f *FSRepository) Close() error {
	return f.file.Close()
}
