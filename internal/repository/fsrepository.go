package repository

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

type storage map[string]string

type FSRepository struct {
	lock sync.Mutex
	urls storage
	file *os.File
}

var _ Repository = (*FSRepository)(nil)

func NewFSRepository(ctx context.Context, file string) (*FSRepository, error) {
	urls, err := readFile(file)
	if err != nil {
		log.Warn().Err(err).Msg("Can not read urls from storage file")
		urls = storage{}
	}

	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	r := &FSRepository{
		lock: sync.Mutex{},
		urls: urls,
		file: f,
	}

	go func() {
		<-ctx.Done()
		if err := r.close(); err != nil {
			log.Error().Err(err).Msg("Can not close file")
		}
	}()
	return r, nil
}

func (f *FSRepository) AddURL(id, url string) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.urls[id] = url
	if err := f.updateFile(); err != nil {
		return err
	}

	return nil
}

func (f *FSRepository) GetURL(id string) (string, error) {
	url, ok := f.urls[id]
	if !ok {
		return "", ErrorNotFound
	}

	return url, nil
}

func (f *FSRepository) close() error {
	return f.file.Close()
}

func (f *FSRepository) updateFile() error {
	if err := f.file.Truncate(0); err != nil {
		return err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}

	enc := json.NewEncoder(f.file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(f.urls); err != nil {
		return err
	}

	return nil
}

func readFile(filename string) (storage, error) {
	file, err := os.Open(filename)
	if err != nil {
		return storage{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	var s storage
	dec := json.NewDecoder(file)
	if err := dec.Decode(&s); err != nil {
		return storage{}, err
	}
	return s, nil
}
