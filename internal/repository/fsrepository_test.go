package repository

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	rmodels "github.com/darrior/urlshortener/internal/repository/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Repository = (*FSRepository)(nil)

func TestFSRepository_AddURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		file *os.File

		urls    urlStorage
		id      string
		url     string
		userID  string
		wantErr bool
		want    urlStorage
	}{
		{
			name: "Add URL to empty map",
			urls: map[string]record{},
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			id:      "test_id",
			url:     "123456",
			userID:  "123",
			wantErr: false,
			want: urlStorage{
				"test_id": {OriginalURL: "123456", UserID: "123"}},
		},
		{
			name: "Add URL to non-empty map",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			urls: urlStorage{
				"id_test": {OriginalURL: "654321", UserID: "123"},
			},
			id:      "test_id",
			url:     "123456",
			userID:  "123",
			wantErr: false,
			want: urlStorage{
				"id_test": {OriginalURL: "654321", UserID: "123"},
				"test_id": {OriginalURL: "123456", UserID: "123"}},
		},
		{
			name: "Overwrite URL in map",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			urls: urlStorage{"test_id": {
				OriginalURL: "654321",
				UserID:      "321",
			}},
			id:      "test_id",
			url:     "123456",
			userID:  "123",
			wantErr: false,
			want: urlStorage{"test_id": {
				OriginalURL: "123456",
				UserID:      "123",
			}},
		},
		{
			name: "Non empty file",
			file: func() *os.File {
				urls := urlStorage{
					"id_test": record{OriginalURL: "654321", UserID: "123"},
				}

				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				enc := json.NewEncoder(f)
				err = enc.Encode(urls)
				assert.NoError(t, err)

				return f
			}(),
			urls:    urlStorage{},
			id:      "test_id",
			url:     "123456",
			userID:  "123",
			wantErr: false,
			want: urlStorage{
				"id_test": {OriginalURL: "654321", UserID: "123"},
				"test_id": {OriginalURL: "123456", UserID: "123"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := os.Remove(tt.file.Name())
				assert.NoError(t, err)
			}()

			f, err := NewFSRepository(tt.file)
			require.NoError(t, err)

			if len(tt.urls) > 0 {
				f.urls = tt.urls
			}

			gotErr := f.AddURL(context.TODO(), tt.userID, tt.id, tt.url)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, f.urls)
		})
	}
}

func TestFSRepository_GetURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		file *os.File
		// Named input parameters for target function.
		urls    urlStorage
		id      string
		want    string
		wantErr bool
	}{
		{
			name: "Valid get from map",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			urls: urlStorage{
				"test": {OriginalURL: "123", UserID: "123"},
			},
			id:      "test",
			want:    "123",
			wantErr: false,
		},
		{
			name: "Non empty file",
			file: func() *os.File {
				urls := urlStorage{
					"test": {OriginalURL: "123", UserID: "123"},
				}

				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				enc := json.NewEncoder(f)
				err = enc.Encode(urls)
				assert.NoError(t, err)

				return f
			}(),
			urls:    urlStorage{},
			id:      "test",
			want:    "123",
			wantErr: false,
		},
		{
			name: "Unexisting id",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			urls:    urlStorage{},
			id:      "test",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFSRepository(tt.file)
			require.NoError(t, err)

			if len(tt.urls) > 0 {
				f.urls = tt.urls
			}

			got, gotErr := f.GetURL(context.TODO(), tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFSRepository_AddURLs(t *testing.T) {
	tests := []struct {
		// Named input parameters for target function.
		name string // description of this test case
		file *os.File

		userID    string
		batchURLs rmodels.BatchURLs
		urls      urlStorage
		wantErr   bool
		want      urlStorage
	}{
		{
			name: "Add URL to empty map",
			urls: map[string]record{},
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			userID: "123",
			batchURLs: rmodels.BatchURLs{
				{
					ID:  "abc",
					URL: "123",
				},
			},
			wantErr: false,
			want: urlStorage{"abc": {
				OriginalURL: "123",
				UserID:      "123",
			}},
		},
		{
			name: "Non empty file",
			file: func() *os.File {
				urls := urlStorage{
					"id_test": {OriginalURL: "654321", UserID: "321"},
				}

				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				enc := json.NewEncoder(f)
				err = enc.Encode(urls)
				assert.NoError(t, err)

				return f
			}(),
			urls:   urlStorage{},
			userID: "123",
			batchURLs: rmodels.BatchURLs{
				{
					ID:  "abc",
					URL: "123",
				},
			},
			wantErr: false,
			want: urlStorage{
				"id_test": {OriginalURL: "654321", UserID: "321"},
				"abc":     {OriginalURL: "123", UserID: "123"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFSRepository(tt.file)
			require.NoError(t, err)

			if len(tt.urls) > 0 {
				f.urls = tt.urls
			}

			gotErr := f.AddURLs(context.Background(), tt.userID, tt.batchURLs)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, f.urls)
		})
	}
}

func TestFSRepository_Count(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		file    *os.File
		want    int
		wantErr bool
	}{
		{
			name: "Empty storage",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			want:    0,
			wantErr: false,
		},
		{
			name: "Non empty storage",
			file: func() *os.File {
				urls := urlStorage{
					"id_test": {OriginalURL: "123", UserID: "123"},
				}

				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				enc := json.NewEncoder(f)
				err = enc.Encode(urls)
				assert.NoError(t, err)

				return f
			}(),
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFSRepository(tt.file)

			require.NoError(t, err)

			got, gotErr := f.Count(context.Background())

			if tt.wantErr {
				require.Error(t, gotErr)
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFSRepository_Close(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		file *os.File
	}{
		{
			name: "Always valid",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFSRepository(tt.file)

			require.NoError(t, err)

			gotErr := f.Close()

			assert.NoError(t, gotErr)
		})
	}
}
