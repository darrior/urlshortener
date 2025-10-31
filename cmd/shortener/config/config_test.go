package config

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_validateListenAddress(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		address string
		want    host
		wantErr bool
	}{
		{
			name:    "Empty string",
			address: "",
			want:    "",
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
			want:    "",
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
			c := Config{
				BaseAddress: url.URL{},
			}
			gotErr := c.validateBaseAddress(tt.address)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, c.BaseAddress.String())
		})
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Flags arguments
		flags []string
		// ENV variables
		env     map[string]string
		want    Config
		wantErr bool
	}{
		{
			name:  "ENV and flags",
			flags: []string{"-a", "127.0.0.1:80", "-b", "http://127.0.0.1:90"},
			env:   map[string]string{"LISTEN_ADDRESS": "127.0.0.1:50", "BASE_ADDRESS": "https://127.0.0.1:9090"},
			want: Config{
				ListenAddress: "127.0.0.1:50",
				BaseAddress: url.URL{
					Scheme: "https",
					Host:   "127.0.0.1:9090",
				},
			},
			wantErr: false,
		},
		{
			name:  "Fallback to flag",
			flags: []string{"-a", "127.0.0.1:80", "-b", "http://127.0.0.1:90"},
			env:   map[string]string{"LISTEN_ADDRESS": "127.0.0.1:50"},
			want: Config{
				ListenAddress: "127.0.0.1:50",
				BaseAddress: url.URL{
					Scheme: "http",
					Host:   "127.0.0.1:90",
				},
			},
			wantErr: false,
		},
		{
			name:  "Ivalid env",
			flags: []string{"-a", "127.0.0.1:80", "-b", "http://127.0.0.1:90"},
			env:   map[string]string{"LISTEN_ADDRESS": "127.0.0.1:1234567", "BASE_ADDRESS": "https://127.0.0.1:9090"},
			want: Config{
				ListenAddress: "127.0.0.1:50",
				BaseAddress: url.URL{
					Scheme: "http",
					Host:   "127.0.0.1:90",
				},
			},
			wantErr: true,
		},
		{
			name:  "Ivalid flags",
			flags: []string{"-a", "127.0.0.1:1234567", "-b", "http://127.0.0.1:90"},
			env:   map[string]string{"LISTEN_ADDRESS": "127.0.0.1:50", "BASE_ADDRESS": "https://127.0.0.1:9090"},
			want: Config{
				ListenAddress: "127.0.0.1:50",
				BaseAddress: url.URL{
					Scheme: "http",
					Host:   "127.0.0.1:90",
				},
			},
			wantErr: true,
		},
		{
			name:  "Fallback to default",
			flags: []string{},
			env:   map[string]string{},
			want: Config{
				ListenAddress: _defaultListenAddress,
				BaseAddress:   _defaultBaseAddress,
			},
			wantErr: false,
		},
		{
			name:  "Flags without env",
			flags: []string{"-a", "127.0.0.1:80", "-b", "http://127.0.0.1:90"},
			env:   map[string]string{},
			want: Config{
				ListenAddress: "127.0.0.1:80",
				BaseAddress: url.URL{
					Scheme: "http",
					Host:   "127.0.0.1:90",
				},
			},
			wantErr: false,
		},
		{
			name:  "Env without flags",
			flags: []string{},
			env:   map[string]string{"LISTEN_ADDRESS": "127.0.0.1:50", "BASE_ADDRESS": "https://127.0.0.1:9090"},
			want: Config{
				ListenAddress: "127.0.0.1:50",
				BaseAddress: url.URL{
					Scheme: "https",
					Host:   "127.0.0.1:9090",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{"test"}
			os.Args = append(os.Args, tt.flags...)

			for k, v := range tt.env {
				err := os.Setenv(k, v)
				assert.NoError(t, err)
			}

			defer func() {
				for k := range tt.env {
					err := os.Unsetenv(k)
					assert.NoError(t, err)
				}
			}()

			got, gotErr := ParseConfig()

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseHost(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h       string
		want    host
		wantErr bool
	}{
		{
			name:    "Valid host with ip",
			h:       "127.0.0.1:8080",
			want:    "127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "Valid host with domain",
			h:       "example.com:8080",
			want:    "example.com:8080",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			h:       "127.0",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid port",
			h:       "127.0.0.1:1234567",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parseHost(tt.h)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
