package repository

import (
	"testing"

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
			gotErr := r.AddURL(tt.id, tt.url)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.url, r.urls[tt.id])
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
			got, gotErr := r.GetURL(tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
