package handler

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size = size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.responseData.status = status
}

var _ http.ResponseWriter = new(loggingResponseWriter)

func logMiddlware(h http.Handler) http.Handler {
	logHandler := func(res http.ResponseWriter, req *http.Request) {
		lres := loggingResponseWriter{
			ResponseWriter: res,
			responseData:   &responseData{},
		}

		start := time.Now()

		h.ServeHTTP(&lres, req)

		duration := time.Since(start)

		log.Info().
			Str("url", req.URL.String()).
			Str("method", req.Method).
			Int64("time in mikros", duration.Microseconds()).
			Msg("request received")
		log.Info().
			Int("status code", lres.responseData.status).
			Int("size", lres.responseData.size).
			Msg("response written")

	}

	return http.HandlerFunc(logHandler)
}

func extractMiddlware(h http.Handler) http.Handler {
	extractHandler := func(res http.ResponseWriter, req *http.Request) {
		encoding := req.Header.Get("content-encoding")
		if encoding == "" {
			h.ServeHTTP(res, req)
			return
		}

		if encoding != "gzip" {
			http.Error(res, "Unsupported content encoding", http.StatusBadRequest)
		}

		r, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Close()

		data, err := io.ReadAll(r)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		req.Body = io.NopCloser(bytes.NewReader(data))
		req.Header.Del("content-encoding")
		req.Header.Set("content-length", strconv.Itoa(len(data)))

		h.ServeHTTP(res, req)
	}

	return http.HandlerFunc(extractHandler)
}
