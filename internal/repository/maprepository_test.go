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
		urls    map[string]string
		id      string
		url     string
		wantErr bool
	}{
		{
			name:    "Add URL to empty map",
			urls:    map[string]string{},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
		},
		{
			name:    "Add URL to non-empty map",
			urls:    map[string]string{"id_test": "654321"},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
		},
		{
			name:    "Overwrite URL in map",
			urls:    map[string]string{"test_id": "654321"},
			id:      "test_id",
			url:     "123456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MapRepository{urls: tt.urls}
			gotErr := r.AddURL(context.TODO(), tt.id, tt.url)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.url, r.urls[tt.id])
		})
	}
}

func TestMapRepository_AddURLs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		batchURLs models.BatchURLs
		want      map[string]string
		wantErr   bool
	}{
		{
			name: "Test AddURLs",
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
			want: map[string]string{
				"abc": "123",
				"cba": "321",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapRepository()
			gotErr := m.AddURLs(context.Background(), tt.batchURLs)

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
		urls    map[string]string
		id      string
		want    string
		wantErr bool
	}{
		{
			name:    "Get url from empty map",
			urls:    map[string]string{},
			id:      "test_id",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Get valid url",
			urls:    map[string]string{"test_id": "123456"},
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
		urls    map[string]string
		want    int
		wantErr bool
	}{
		{
			name:    "Empty map",
			urls:    map[string]string{},
			want:    0,
			wantErr: false,
		},
		{
			name: "Non empty map",
			urls: map[string]string{
				"abc": "123",
				"cba": "321",
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
