package repository

import (
	"context"
	"testing"

	"github.com/darrior/urlshortener/internal/models"
	"github.com/stretchr/testify/assert"
)

var _ Repository = (*MapRepository)(nil)

func TestMapRepository_AddURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		urls    urlStorage
		userID  string
		id      string
		url     string
		wantErr bool
	}{
		{
			name:    "Add URL to empty map",
			urls:    urlStorage{},
			userID:  "123",
			id:      "test_id",
			url:     "123456",
			wantErr: false,
		},
		{
			name: "Add URL to non-empty map",
			urls: urlStorage{
				"id_test": {OriginalURL: "654321", UserID: "321"},
			},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
		},
		{
			name: "Overwrite URL in map",
			urls: urlStorage{
				"test_id": {OriginalURL: "654321", UserID: "321"},
			},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MapRepository{urls: tt.urls}
			gotErr := r.AddURL(context.TODO(), tt.userID, tt.id, tt.url)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.url, r.urls[tt.id].OriginalURL)
		})
	}
}

func TestMapRepository_AddURLs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID    string
		batchURLs models.BatchURLs
		want      urlStorage
		wantErr   bool
	}{
		{
			name:   "Test AddURLs",
			userID: "123",
			batchURLs: models.BatchURLs{
				{
					ID:  "abc",
					URL: "123",
				},
				{
					ID:  "cba",
					URL: "321",
				},
			},
			want: urlStorage{
				"abc": {OriginalURL: "123", UserID: "123"},
				"cba": {OriginalURL: "321", UserID: "123"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapRepository()
			gotErr := m.AddURLs(context.Background(), tt.userID, tt.batchURLs)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, m.urls)
		})
	}
}

func TestMapRepository_GetURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		urls    urlStorage
		id      string
		want    string
		wantErr bool
	}{
		{
			name:    "Get url from empty map",
			urls:    urlStorage{},
			id:      "test_id",
			want:    "",
			wantErr: true,
		},
		{
			name: "Get valid url",
			urls: urlStorage{
				"test_id": {OriginalURL: "123456", UserID: "123"},
			},
			id:      "test_id",
			want:    "123456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MapRepository{urls: tt.urls}
			got, gotErr := r.GetURL(context.TODO(), tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapRepository_Count(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		urls    urlStorage
		want    int
		wantErr bool
	}{
		{
			name:    "Empty map",
			urls:    urlStorage{},
			want:    0,
			wantErr: false,
		},
		{
			name: "Non empty map",
			urls: urlStorage{
				"abc": {OriginalURL: "123", UserID: "123"},
				"cba": {OriginalURL: "321", UserID: "123"},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MapRepository{
				urls: tt.urls,
			}
			got, gotErr := m.Count(context.Background())

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapRepository_Ping(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			"Always valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapRepository()
			gotErr := m.Ping(context.Background())
			assert.NoError(t, gotErr)
		})
	}
}

func TestMapRepository_Close(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{"Always valid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapRepository()
			gotErr := m.Close()

			assert.NoError(t, gotErr)
		})
	}
}
