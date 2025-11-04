package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
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

var _ http.ResponseWriter = new(loggingResponseWriter)

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size = size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.ResponseWriter.WriteHeader(status)
	l.responseData.status = status
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

var _ http.ResponseWriter = new(gzipResponseWriter)

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

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
		if !strings.Contains(req.Header.Get("content-encoding"), "gzip") {
			h.ServeHTTP(res, req)
			return
		}

		b := req.Body
		defer func() {
			_ = b.Close()
		}()

		r, err := gzip.NewReader(b)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = r.Close()
		}()

		req.Body = r

		h.ServeHTTP(res, req)
	}

	return http.HandlerFunc(extractHandler)
}

func compressMiddlware(h http.Handler) http.Handler {
	compressHandler := func(res http.ResponseWriter, req *http.Request) {
		if !checkEncoding(req.Header.Values("accept-encoding"), "gzip") {
			h.ServeHTTP(res, req)
			return
		}

		var data []byte
		buf := bytes.NewBuffer(data)
		h.ServeHTTP(&gzipResponseWriter{res, buf}, req)

		if !strings.Contains(res.Header().Get("content-type"), "application/json") &&
			!strings.Contains(res.Header().Get("content-type"), "text/html") {
			_, _ = res.Write(data)
			return
		}

		w, err := gzip.NewWriterLevel(res, gzip.BestSpeed)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		n, err := w.Write(data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := w.Close(); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("conent-length", strconv.Itoa(n))
		res.Header().Set("content-encoding", "gzip")
	}

	return http.HandlerFunc(compressHandler)
}

func checkEncoding(encs []string, enc string) bool {
	for _, e := range encs {
		if strings.Contains(e, enc) {
			return true
		}
	}

	return false
}
