package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		file    *os.File
		data    any
		want    any
		wantErr bool
	}{
		{
			name: "Read map",
			file: func() *os.File {
				f, err := os.CreateTemp("", "test-*.txt")
				assert.NoError(t, err)

				e := json.NewEncoder(f)
				err = e.Encode(map[string]string{"Hello": "world"})
				assert.NoError(t, err)

				return f
			}(),
			data:    map[string]string{},
			want:    map[string]any{"Hello": "world"},
			wantErr: false,
		},
		{
			name:    "All nil",
			file:    nil,
			data:    nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.file == nil {
					return
				}
				err := os.Remove(tt.file.Name())
				assert.NoError(t, err)
			}()
			gotErr := ReadFile(tt.file, &tt.data)
			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, tt.data)
		})
	}
}

func TestUpdateFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		file    *os.File
		data    any
		want    []byte
		wantErr bool
	}{
		{
			name: "Write map",
			file: func() *os.File {
				f, err := os.CreateTemp("", "test-*.txt")
				assert.NoError(t, err)
				return f
			}(),
			data:    map[string]string{"Hello": "World"},
			want:    []byte("{\n  \"Hello\": \"World\"\n}\n"),
			wantErr: false,
		},
		{
			name:    "All nil",
			file:    nil,
			data:    nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.file == nil {
					return
				}

				err := os.Remove(tt.file.Name())
				assert.NoError(t, err)
			}()
			gotErr := UpdateFile(tt.file, tt.data)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			err := tt.file.Close()
			assert.NoError(t, err)
			data, err := os.ReadFile(tt.file.Name())
			assert.NoError(t, err)
			assert.Equal(t, tt.want, data)
		})
	}
}
