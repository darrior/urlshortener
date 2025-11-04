package repository

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSRepository_AddURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		file string

		urls    storage
		id      string
		url     string
		wantErr bool
		want    storage
	}{
		{
			name: "Add URL to empty map",
			urls: map[string]string{},
			file: func() string {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				err = f.Close()
				assert.NoError(t, err)

				return f.Name()
			}(),
			id:      "test_id",
			url:     "123456",
			wantErr: false,
			want:    storage{"test_id": "123456"},
		},
		{
			name: "Add URL to non-empty map",
			file: func() string {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				err = f.Close()
				assert.NoError(t, err)

				return f.Name()
			}(),
			urls:    storage{"id_test": "654321"},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
			want:    storage{"id_test": "654321", "test_id": "123456"},
		},
		{
			name: "Overwrite URL in map",
			file: func() string {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				err = f.Close()
				assert.NoError(t, err)

				return f.Name()
			}(),
			urls:    storage{"test_id": "654321"},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
			want:    storage{"test_id": "123456"},
		},
		{
			name: "Non empty file",
			file: func() string {
				urls := storage{"id_test": "654321"}

				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				enc := json.NewEncoder(f)
				err = enc.Encode(urls)
				assert.NoError(t, err)

				err = f.Close()
				assert.NoError(t, err)

				return f.Name()
			}(),
			urls:    storage{},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
			want:    storage{"id_test": "654321", "test_id": "123456"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := os.Remove(tt.file)
				assert.NoError(t, err)
			}()

			f, err := NewFSRepository(context.Background(), tt.file)
			assert.NoError(t, err)

			if len(tt.urls) > 0 {
				f.urls = tt.urls
			}

			gotErr := f.AddURL(tt.id, tt.url)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, f.urls)
			assert.FileExists(t, tt.file)
		})
	}
}

func TestFSRepository_GetURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		file string
		// Named input parameters for target function.
		urls    storage
		id      string
		want    string
		wantErr bool
	}{
		{
			name: "Valid get from map",
			file: func() string {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				err = f.Close()
				assert.NoError(t, err)

				return f.Name()
			}(),
			urls:    storage{"test": "123"},
			id:      "test",
			want:    "123",
			wantErr: false,
		},
		{
			name: "Non empty file",
			file: func() string {
				urls := storage{"test": "123"}

				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				enc := json.NewEncoder(f)
				err = enc.Encode(urls)
				assert.NoError(t, err)

				err = f.Close()
				assert.NoError(t, err)

				return f.Name()
			}(),
			urls:    storage{},
			id:      "test",
			want:    "123",
			wantErr: false,
		},
		{
			name: "Unexisting id",
			file: func() string {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				err = f.Close()
				assert.NoError(t, err)

				return f.Name()
			}(),
			urls:    storage{},
			id:      "test",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFSRepository(context.Background(), tt.file)
			assert.NoError(t, err)

			if len(tt.urls) > 0 {
				f.urls = tt.urls
			}

			got, gotErr := f.GetURL(tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
