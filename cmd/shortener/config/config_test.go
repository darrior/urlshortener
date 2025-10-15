package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_validateListenAddress(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		address string
		want    string
		wantErr bool
	}{
		{
			name:    "Empty string",
			address: "",
			want:    _defaultListenAddress,
			wantErr: false,
		},
		{
			name:    "Valid address",
			address: "127.0.0.1:8080",
			want:    "127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "Valid address with host",
			address: "example.com:80",
			want:    "example.com:80",
			wantErr: false,
		},
		{
			name:    "Invalid address",
			address: "http://127.0.0.1:8080",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{}
			gotErr := c.validateListenAddress(tt.address)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, c.ListenAddress)
		})
	}
}

func TestConfig_validateBaseAddress(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		address string
		want    string
		wantErr bool
	}{
		{
			name:    "Empty string",
			address: "",
			want:    _defaultBaseAddress,
			wantErr: false,
		},
		{
			name:    "Valid address",
			address: "http://127.0.0.1:8080",
			want:    "http://127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "Valid address with domain name",
			address: "http://example.com",
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "Valid address with path",
			address: "http://127.0.0.1:8080/test/path",
			want:    "http://127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "Valid address with query",
			address: "http://127.0.0.1:8080/?hello=1",
			want:    "http://127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "Valid address without scheme",
			address: "//127.0.0.1:8080/?hello=1",
			want:    "//127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "Invalid address",
			address: "127.0.0.1:8080/?hello=1",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{}
			gotErr := c.validateBaseAddress(tt.address)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, c.BaseAddress)
		})
	}
}
