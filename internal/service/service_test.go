package service

import (
	"context"
	"testing"

	"github.com/darrior/urlshortener/internal/models"
	"github.com/darrior/urlshortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

var _ IService = (*Service)(nil)

type testRepository struct {
	urls map[string]string
}

var _ repository.Repository = (*testRepository)(nil)

func (t *testRepository) AddURL(_ context.Context, id, url string) error {
	t.urls[id] = url

	return nil
}

func (t *testRepository) GetURL(_ context.Context, id string) (string, error) {
	url, ok := t.urls[id]
	if !ok {
		return "", repository.ErrorNotFound
	}
	return url, nil
}

func (t *testRepository) AddURLs(_ context.Context, batchURLs models.BatchURLs) error {
	for _, url := range batchURLs {
		t.urls[url.ID] = url.URL
	}

	return nil
}

func (t *testRepository) Count(_ context.Context) (int, error) {
	return len(t.urls), nil
}

func (t *testRepository) Ping(_ context.Context) error {
	return nil
}

func (t *testRepository) Close() error {
	return nil
}

func TestService_AddURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		data repository.Repository
		// Named input parameters for target function.
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "Add url",
			data:    &testRepository{make(map[string]string)},
			url:     "http://example.com",
			want:    "http://127.0.0.1:8080/AAAAAAA",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.data, "http://127.0.0.1:8080")
			got, gotErr := s.AddURL(context.TODO(), tt.url)

			if tt.wantErr {
				assert.EqualError(t, ErrorCannotAddURL, gotErr.Error())
				return
			}

			assert.NoError(t, gotErr)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_AddURLs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		data        repository.Repository
		baseAddress string
		// Named input parameters for target function.
		longURLs models.ShortenerBatchRequest
		want     models.ShortenerBatchResponse
		wantErr  bool
	}{
		{
			name:        "Empty repository",
			data:        &testRepository{make(map[string]string)},
			baseAddress: "http://127.0.0.1:8080",
			longURLs: models.ShortenerBatchRequest{
				{
					CorrelationID: "a",
					OriginalURL:   "https://example.com",
				},
				{
					CorrelationID: "b",
					OriginalURL:   "http://example.com",
				},
			},
			want: models.ShortenerBatchResponse{
				{
					CorrelationID: "a",
					ShortURL:      "http://127.0.0.1:8080/AAAAAAA",
				},
				{
					CorrelationID: "b",
					ShortURL:      "http://127.0.0.1:8080/AAAAAAB",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.data, tt.baseAddress)
			got, gotErr := s.AddURLs(context.Background(), tt.longURLs)

			if tt.wantErr {
				assert.EqualError(t, ErrorCannotAddURL, gotErr.Error())
				return
			}

			assert.NoError(t, gotErr)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		data repository.Repository
		// Named input parameters for target function.
		id      string
		want    string
		wantErr bool
	}{
		{
			name:    "Get url from empty repository",
			data:    &testRepository{urls: map[string]string{}},
			id:      "test_id",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Get existing url",
			data:    &testRepository{urls: map[string]string{"AAAAAAA": "http://example.com"}},
			id:      "AAAAAAA",
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "Get unexisting url",
			data:    &testRepository{urls: map[string]string{"AAAAAAA": "http://example.com"}},
			id:      "BBBBBBB",
			want:    "http://example.com",
			wantErr: true,
		},
		{
			name: "Get existing url",
			data: &testRepository{urls: map[string]string{
				"AAAAAAA": "http://example.com",
				"BBBBBBB": "http://example1.com",
				"CCCCCCC": "http://example2.com",
			}},
			id:      "BBBBBBB",
			want:    "http://example1.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.data, "http://127.0.0.1:8080")
			got, gotErr := s.GetURL(context.TODO(), tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
