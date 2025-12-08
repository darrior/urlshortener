package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

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
	status int
	writer io.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipResponseWriter) WriteHeader(status int) {
	g.status = status
}

func checkEncoding(encs []string, enc string) bool {
	for _, e := range encs {
		if strings.Contains(e, enc) {
			return true
		}
	}

	return false
}

func (h *handler) checkAuthCookies(cookies []*http.Cookie) (*http.Cookie, error) {
	var errs []error
	for _, cookie := range cookies {
		if valid, err := h.service.ValidateToken(cookie.Value); err != nil && valid {
			return cookie, nil
		} else if err != nil {
			errs = append(errs, err)
		}
	}

	return nil, fmt.Errorf("errors in token validation: %w", errors.Join(errs...))
}
