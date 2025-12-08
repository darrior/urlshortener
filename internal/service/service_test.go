package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/darrior/urlshortener/internal/mocks"
	"github.com/darrior/urlshortener/internal/models/api"
	"github.com/darrior/urlshortener/internal/repository"
	rmodels "github.com/darrior/urlshortener/internal/repository/models"
)

var _ IService = new(Service)

func TestService_AddURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID  string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "Add url",
			userID:  "123",
			url:     "http://example.com",
			want:    "http://127.0.0.1:8080/AAAAAAA",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		m.EXPECT().
			AddURL(gomock.Any(), "123", "AAAAAAA", "http://example.com").
			Return(nil)
		m.EXPECT().Count(gomock.Any()).Return(0, nil)

		t.Run(tt.name, func(t *testing.T) {
			s := NewService(m, "http://127.0.0.1:8080", "123")
			got, gotErr := s.AddURL(context.TODO(), tt.userID, tt.url)

			if tt.wantErr {
				assert.EqualError(t, ErrorCannotAddURL, gotErr.Error())
				return
			}

			assert.NoError(t, gotErr)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_AddURLs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		baseAddress string
		// Named input parameters for target function.
		userID   string
		longURLs api.ShortenerBatchRequest
		want     api.ShortenerBatchResponse
		wantErr  bool
	}{
		{
			name:        "Empty repository",
			baseAddress: "http://127.0.0.1:8080",
			userID:      "123",
			longURLs: api.ShortenerBatchRequest{
				{
					CorrelationID: "a",
					OriginalURL:   "https://example.com",
				},
				{
					CorrelationID: "b",
					OriginalURL:   "http://example.com",
				},
			},
			want: api.ShortenerBatchResponse{
				{
					CorrelationID: "a",
					ShortURL:      "http://127.0.0.1:8080/AAAAAAA",
				},
				{
					CorrelationID: "b",
					ShortURL:      "http://127.0.0.1:8080/AAAAAAB",
				},
			},
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().AddURLs(
		gomock.Any(),
		"123",
		rmodels.BatchURLs{
			{
				ID:  "AAAAAAA",
				URL: "https://example.com",
			},
			{
				ID:  "AAAAAAB",
				URL: "http://example.com",
			},
		}).
		Return(
			nil,
		)
	count := 0
	m.EXPECT().Count(gomock.Any()).DoAndReturn(func(_ any) (int, error) { count += 1; return count - 1, nil })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(m, tt.baseAddress, "")
			got, gotErr := s.AddURLs(context.Background(), tt.userID, tt.longURLs)

			if tt.wantErr {
				assert.EqualError(t, ErrorCannotAddURL, gotErr.Error())
				return
			}

			assert.NoError(t, gotErr)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      string
		want    string
		wantErr bool
	}{
		{
			name:    "Get url from empty repository",
			id:      "test_id",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Get existing url",
			id:      "AAAAAAA",
			want:    "http://example.com",
			wantErr: false,
		},
		{
			name:    "Get unexisting url",
			id:      "BBBBBBB",
			want:    "http://example.com",
			wantErr: true,
		},
		{
			name:    "Get existing url",
			id:      "CCCCCCC",
			want:    "http://example1.com",
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	m.EXPECT().GetURL(context.TODO(), "test_id").Return("", repository.ErrorNotFound)
	m.EXPECT().GetURL(context.TODO(), "AAAAAAA").Return("http://example.com", nil)
	m.EXPECT().GetURL(context.TODO(), "BBBBBBB").Return("", repository.ErrorNotFound)
	m.EXPECT().GetURL(context.TODO(), "CCCCCCC").Return("http://example1.com", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(m, "http://127.0.0.1:8080", "")
			got, gotErr := s.GetURL(context.TODO(), tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
