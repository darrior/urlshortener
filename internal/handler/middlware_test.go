package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Checks for interface implementations.
var _ http.ResponseWriter = new(loggingResponseWriter)

var _ http.ResponseWriter = new(gzipResponseWriter)

type mwant struct {
	status int
	data   []byte
}

func echoHandler(res http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := req.Body.Close(); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.Header().Set("content-length", strconv.Itoa(len(data)))
	_, _ = res.Write(data)
	res.WriteHeader(http.StatusOK)
}

func Test_extractMiddlware(t *testing.T) {

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    http.Handler
		req  *http.Request
		want mwant
	}{
		{
			name: "Valid request",
			h:    http.HandlerFunc(echoHandler),
			req: func() *http.Request {
				data := []byte("Hello, world!")

				var body bytes.Buffer

				g, _ := gzip.NewWriterLevel(&body, gzip.BestSpeed)
				_, err := g.Write(data)
				assert.NoError(t, err)

				err = g.Close()
				assert.NoError(t, err)

				data, err = io.ReadAll(&body)
				assert.NoError(t, err)

				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
				r.Header.Set("content-encoding", "gzip")
				return r
			}(),
			want: mwant{
				status: http.StatusOK,
				data:   []byte("Hello, world!"),
			},
		},
		{
			name: "Invalid gzip",
			h:    http.HandlerFunc(echoHandler),
			req: func() *http.Request {
				data := []byte("Hello, world!")

				var body bytes.Buffer

				g, _ := gzip.NewWriterLevel(&body, gzip.BestSpeed)
				_, err := g.Write(data)
				assert.NoError(t, err)

				data, err = io.ReadAll(&body)
				assert.NoError(t, err)

				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
				r.Header.Set("content-encoding", "gzip")
				return r
			}(),
			want: mwant{
				status: http.StatusInternalServerError,
				data:   []byte("unexpected EOF\n"),
			},
		},
		{
			name: "Invalid gzip header",
			h:    http.HandlerFunc(echoHandler),
			req: func() *http.Request {
				data := []byte("Hello, world!")

				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
				r.Header.Set("content-encoding", "gzip")
				return r
			}(),
			want: mwant{
				status: http.StatusInternalServerError,
				data:   []byte("gzip: invalid header\n"),
			},
		},
		{
			name: "Request without encoding",
			h:    http.HandlerFunc(echoHandler),
			req: func() *http.Request {
				data := []byte("Hello, world!")

				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
				return r
			}(),
			want: mwant{
				status: http.StatusOK,
				data:   []byte("Hello, world!"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()

			gotHandler := extractMiddlware(tt.h)
			gotHandler.ServeHTTP(res, tt.req)

			resp := res.Result()
			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			err = resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.want.status, resp.StatusCode)
			assert.Equal(t, tt.want.data, data)
		})
	}
}

func Test_compressMiddlware(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    http.Handler
		req  *http.Request
		want mwant
	}{
		{
			name: "Valid request",
			h:    http.HandlerFunc(echoHandler),
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("Hello, world!")))
				r.Header.Set("Accept-Encoding", "gzip")
				return r

			}(),
			want: func() mwant {
				w := mwant{}
				w.status = http.StatusOK

				var b bytes.Buffer
				g, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
				assert.NoError(t, err)

				_, err = g.Write([]byte("Hello, world!"))
				assert.NoError(t, err)

				err = g.Close()
				assert.NoError(t, err)

				w.data = b.Bytes()

				return w
			}(),
		},
		{
			name: "No accept",
			h:    http.HandlerFunc(echoHandler),
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("Hello, world!")))
				return r

			}(),
			want: mwant{
				status: http.StatusOK,
				data:   []byte("Hello, world!"),
			},
		},
		{
			name: "Wrong content-type",
			h: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				data, err := io.ReadAll(req.Body)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}

				res.Header().Set("content-type", "text/plain")
				res.Header().Set("content-length", strconv.Itoa(len(data)))
				_, _ = res.Write(data)
				res.WriteHeader(http.StatusOK)
			}),
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("Hello, world!")))
				r.Header.Set("Accept-Encoding", "gzip")
				return r

			}(),
			want: mwant{
				status: http.StatusOK,
				data:   []byte("Hello, world!"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()

			got := compressMiddlware(tt.h)
			got.ServeHTTP(res, tt.req)

			resp := res.Result()
			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			err = resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.want.status, resp.StatusCode)
			assert.Equal(t, tt.want.data, data)
		})
	}
}

func Test_checkEncoding(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		encs []string
		enc  string
		want bool
	}{
		{
			name: "Contains gzip",
			encs: []string{"gzip"},
			enc:  "gzip",
			want: true,
		},
		{
			name: "Does not contains gzip",
			encs: []string{"zip"},
			enc:  "gzip",
			want: false,
		},
		{
			name: "Multiple values",
			encs: []string{"zip", "deflate", "gzip"},
			enc:  "gzip",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkEncoding(tt.encs, tt.enc)

			assert.Equal(t, tt.want, got)
		})
	}
}
