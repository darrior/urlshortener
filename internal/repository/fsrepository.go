package repository

import (
	"context"
	"os"
	"sync"

	"github.com/darrior/urlshortener/internal/models"
	"github.com/darrior/urlshortener/internal/repository/storage"
	"github.com/rs/zerolog/log"
)

type FSRepository struct {
	MapRepository
	file *os.File
}

func NewFSRepository(file *os.File) (*FSRepository, error) {
	var urls urlStorage
	if err := storage.ReadFile(file, &urls); err != nil {
		log.Warn().Err(err).Msg("Can not read urls from storage file")
		urls = urlStorage{}
	}

	r := &FSRepository{
		MapRepository: MapRepository{
			lock: sync.Mutex{},
			urls: urls,
		},
		file: file,
	}

	return r, nil
}

func (f *FSRepository) AddURL(_ context.Context, userID, id, url string) error {
	if err := f.MapRepository.AddURL(context.Background(), userID, id, url); err != nil {
		return err
	}

	if err := storage.UpdateFile(f.file, f.urls); err != nil {
		return err
	}

	return nil
}

func (f *FSRepository) AddURLs(_ context.Context, userID string, batchURLs models.BatchURLs) error {
	if err := f.MapRepository.AddURLs(context.TODO(), userID, batchURLs); err != nil {
		return err
	}

	if err := storage.UpdateFile(f.file, f.urls); err != nil {
		return err
	}

	return nil
}

func (f *FSRepository) Close() error {
	return f.file.Close()
}
