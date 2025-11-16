package repository

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSRepository_AddURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		file *os.File

		urls    urlStorage
		id      string
		url     string
		wantErr bool
		want    urlStorage
	}{
		{
			name: "Add URL to empty map",
			urls: map[string]string{},
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			id:      "test_id",
			url:     "123456",
			wantErr: false,
			want:    urlStorage{"test_id": "123456"},
		},
		{
			name: "Add URL to non-empty map",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			urls:    urlStorage{"id_test": "654321"},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
			want:    urlStorage{"id_test": "654321", "test_id": "123456"},
		},
		{
			name: "Overwrite URL in map",
			file: func() *os.File {
				f, err := os.CreateTemp("", "data-*.json")
				assert.NoError(t, err)

				return f
			}(),
			urls:    urlStorage{"test_id": "654321"},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
			want:    urlStorage{"test_id": "123456"},
		},
		{
			name: "Non empty file",
			file: func() *os.File {
				urls := urlStorage{"id_test": "654321"}

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
			wantErr: false,
			want:    urlStorage{"id_test": "654321", "test_id": "123456"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := os.Remove(tt.file.Name())
				assert.NoError(t, err)
			}()

			f, err := NewFSRepository(tt.file)
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
			urls:    urlStorage{"test": "123"},
			id:      "test",
			want:    "123",
			wantErr: false,
		},
		{
			name: "Non empty file",
			file: func() *os.File {
				urls := urlStorage{"test": "123"}

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
