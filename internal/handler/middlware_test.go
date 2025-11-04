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

func emptyHandler(res http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.Header().Set("content-length", strconv.Itoa(len(data)))
	_, _ = res.Write(data)
	res.WriteHeader(http.StatusOK)
}

func Test_extractMiddlware(t *testing.T) {
	type want struct {
		status int
		data   []byte
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		h    http.Handler
		req  *http.Request
		want want
	}{
		{
			name: "Valid request",
			h:    http.HandlerFunc(emptyHandler),
			req: func() *http.Request {
				data := []byte("Hello, world!")
				var body []byte
				wbody := bytes.NewBuffer(body)
				g, _ := gzip.NewWriterLevel(wbody, gzip.BestSpeed)
				_, err := g.Write(data)
				assert.NoError(t, err)

				r := httptest.NewRequest(http.MethodPost, "/", wbody)
				return r
			}(),
			want: want{
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
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compressMiddlware(tt.h)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("compressMiddlware() = %v, want %v", got, tt.want)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkEncoding(tt.encs, tt.enc)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("checkEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}
