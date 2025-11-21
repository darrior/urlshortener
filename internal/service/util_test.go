package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateURLID(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		number int
		want   string
	}{
		{
			name:   "Zero",
			number: 0,
			want:   "AAAAAAA",
		},
		{
			name:   "Fifty Two",
			number: 51,
			want:   "AAAAAAz",
		},
		{
			name:   "Sixty two",
			number: 62,
			want:   "AAAAABA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateURLID(tt.number)

			assert.Equal(t, tt.want, got)
		})
	}
}
