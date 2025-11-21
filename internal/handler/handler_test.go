package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/darrior/urlshortener/internal/models"
	"github.com/darrior/urlshortener/internal/service"
	"github.com/stretchr/testify/assert"
)

type testService struct {
	urls map[string]string
	ping bool
}

var _ service.IService = (*testService)(nil)

func (t *testService) AddURL(_ context.Context, id string) (string, error) {
	return "http://127.0.0.1:8080/AAAAAAA", nil
}

func (t *testService) AddURLs(ctx context.Context, longURLs models.ShortenerBatchRequest) (shortURLs models.ShortenerBatchResponse, err error) {
	panic("unimplemented")
}

func (t *testService) GetURL(_ context.Context, id string) (string, error) {
	url, ok := t.urls[id]
	if !ok {
		return "", errors.New("")
	}

	return url, nil
}

func (t *testService) Ping(ctx context.Context) error {
	if !t.ping {
		return errors.New("")
	}
	return nil
}

type hwant struct {
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
		want hwant
	}{
		{
			name: "GET request to root",
			h:    handler{service: &testService{}},
			req:  httptest.NewRequest(http.MethodGet, "/", nil),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "DELETE request to root",
			h:    handler{service: &testService{}},
			req:  httptest.NewRequest(http.MethodDelete, "/", nil),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "GET request to some path",
			h:    handler{service: &testService{}},
			req:  httptest.NewRequest(http.MethodGet, "/first/second", nil),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Invalid request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			tt.h.errorHandler(res, tt.req)

			assert.Equal(t, tt.want.status, res.Code)
			assert.Equal(t, tt.want.data, res.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header().Get("content-type"))

		})
	}
}

func Test_handler_postURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    handler
		req  *http.Request
		want hwant
	}{
		{
			name: "Empty POST request",
			h:    handler{service: &testService{}},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				r.Header.Set("content-type", "test")
				return r
			}(),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Content type must be \"text/plain\"\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty POST request with header",
			h:    handler{service: &testService{}},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				r.Header.Add("content-type", "text/plain; charset=utf-8")
				return r
			}(),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Invalid URL\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Valid request",
			h:    handler{service: &testService{urls: make(map[string]string)}},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/", strings.NewReader("https://example.com"))
				r.Header.Add("content-type", "text/plain")

				return r
			}(),
			want: hwant{
				status:      http.StatusCreated,
				data:        "http://127.0.0.1:8080/AAAAAAA",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()

			tt.h.postURL(res, tt.req)

			assert.Equal(t, tt.want.status, res.Code)
			assert.Equal(t, tt.want.data, res.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header().Get("content-type"))
		})
	}
}

func Test_handler_postAPIShorten(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    handler
		req  *http.Request
		want hwant
	}{

		{
			name: "Empty POST request",
			h:    handler{service: &testService{}},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				r.Header.Set("content-type", "test")
				return r
			}(),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Content type must be \"application/json\"\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty POST request with header",
			h:    handler{service: &testService{}},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				r.Header.Add("content-type", "application/json")
				return r
			}(),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Can not unmarshal JSON\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Invalid URL",
			h:    handler{service: &testService{}},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"url": ""}`))
				r.Header.Add("content-type", "application/json")
				return r
			}(),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Invalid URL\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Valid request",
			h:    handler{service: &testService{urls: make(map[string]string)}},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/", strings.NewReader(`{"url": "http://example.com"}`))
				r.Header.Add("content-type", "application/json")

				return r
			}(),
			want: hwant{
				status:      http.StatusCreated,
				data:        `{"result":"http://127.0.0.1:8080/AAAAAAA"}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()

			tt.h.postAPIShorten(res, tt.req)

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
		want hwant
	}{
		{
			name: "Unkonwn short URL",
			h:    handler{service: &testService{urls: make(map[string]string)}},
			req:  httptest.NewRequest(http.MethodGet, "/", nil),
			want: hwant{
				status:      http.StatusBadRequest,
				data:        "Short URL not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Valid request",
			h:    handler{service: &testService{urls: map[string]string{"AAAAAAA": "https://example.com"}}},
			req:  validReq,
			want: hwant{
				status:      http.StatusTemporaryRedirect,
				data:        "<a href=\"https://example.com\">Temporary Redirect</a>.\n\n",
				contentType: "text/html; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()

			tt.h.getFullURL(res, tt.req)

			assert.Equal(t, tt.want.status, res.Code)
			assert.Equal(t, tt.want.data, res.Body.String())
			assert.Equal(t, tt.want.contentType, res.Header().Get("content-type"))
		})
	}
}

func Test_handler_getPing(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    handler
		req  *http.Request
		want hwant
	}{
		{
			name: "Ping OK",
			h:    handler{service: &testService{ping: true}},
			req:  httptest.NewRequest(http.MethodGet, "/ping", nil),
			want: hwant{
				status: http.StatusOK,
			},
		},
		{
			name: "Ping not OK",
			h:    handler{service: &testService{ping: false}},
			req:  httptest.NewRequest(http.MethodGet, "/ping", nil),
			want: hwant{
				status: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			tt.h.getPing(res, tt.req)

			assert.Equal(t, tt.want.status, res.Code)
		})
	}
}
