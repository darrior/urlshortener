package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/darrior/urlshortener/internal/service"
	"github.com/stretchr/testify/assert"
)

type testService struct {
	urls map[string]string
}

var _ service.IService = (*testService)(nil)

func (t *testService) AddURL(id string) (string, error) {
	return "AAAAAAA", nil
}

func (t *testService) GetURL(id string) (string, error) {
	url, ok := t.urls[id]
	if !ok {
		return "", errors.New("")
	}

	return url, nil
}

type want struct {
	status      int
	data        string
	contentType string
}

func Test_handler_errorHandler(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    handler
		req  *http.Request
		want want
	}{
		{
			name: "GET request to root",
			h:    handler{service: &testService{}},
			req:  httptest.NewRequest(http.MethodGet, "/", nil),
			want: want{
				status:      http.StatusBadRequest,
				data:        "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "DELETE request to root",
			h:    handler{service: &testService{}},
			req:  httptest.NewRequest(http.MethodDelete, "/", nil),
			want: want{
				status:      http.StatusBadRequest,
				data:        "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "GET request to some path",
			h:    handler{service: &testService{}},
			req:  httptest.NewRequest(http.MethodGet, "/first/second", nil),
			want: want{
				status:      http.StatusBadRequest,
				data:        "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			res := httptest.NewRecorder()
			tt.h.errorHandler(res, tt.req)

			assert.Equal(t, tt.want.status, res.Code)
			assert.Equal(t, tt.want.data, res.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header().Get("content-type"))

		})
	}
}

func Test_handler_postURL(t *testing.T) {
	emptyReq := httptest.NewRequest(http.MethodPost, "/", nil)
	emptyReq.Header.Add("content-type", "text/plain")

	validReq := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/", strings.NewReader("https://example.com"))
	validReq.Header.Add("content-type", "text/plain")

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    handler
		req  *http.Request
		want want
	}{
		{
			name: "Empty POST request",
			h:    handler{service: &testService{}},
			req:  httptest.NewRequest(http.MethodPost, "/", nil),
			want: want{
				status:      http.StatusBadRequest,
				data:        "Content type must be \"text/plain\"\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty POST request with header",
			h:    handler{service: &testService{}},
			req:  emptyReq,
			want: want{
				status:      http.StatusBadRequest,
				data:        "Invalid URL\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Valid request",
			h:    handler{service: &testService{urls: make(map[string]string)}},
			req:  validReq,
			want: want{
				status:      http.StatusCreated,
				data:        "http://127.0.0.1:8080/AAAAAAA",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			res := httptest.NewRecorder()

			tt.h.postURL(res, tt.req)

			assert.Equal(t, tt.want.status, res.Code)
			assert.Equal(t, tt.want.data, res.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header().Get("content-type"))
		})
	}
}

func Test_handler_getFullURL(t *testing.T) {
	validReq := httptest.NewRequest(http.MethodGet, "/", nil)
	validReq.SetPathValue("url_id", "AAAAAAA")

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    handler
		req  *http.Request
		want want
	}{
		{
			name: "Unkonwn short URL",
			h:    handler{service: &testService{urls: make(map[string]string)}},
			req:  httptest.NewRequest(http.MethodGet, "/", nil),
			want: want{
				status:      http.StatusBadRequest,
				data:        "Short URL not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Valid request",
			h:    handler{service: &testService{urls: map[string]string{"AAAAAAA": "https://example.com"}}},
			req:  validReq,
			want: want{
				status:      http.StatusTemporaryRedirect,
				data:        "<a href=\"https://example.com\">Temporary Redirect</a>.\n\n",
				contentType: "text/html; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			res := httptest.NewRecorder()

			tt.h.getFullURL(res, tt.req)

			assert.Equal(t, tt.want.status, res.Code)
			assert.Equal(t, tt.want.data, res.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header().Get("content-type"))
		})
	}
}
