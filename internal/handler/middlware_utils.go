package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	smodels "github.com/darrior/urlshortener/internal/service/models"
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

type fakeResponseWriter struct {
	http.ResponseWriter
	status int
	writer io.Writer
}

func (f *fakeResponseWriter) Write(data []byte) (int, error) {
	return f.writer.Write(data)
}

func (f *fakeResponseWriter) WriteHeader(status int) {
	f.status = status
}

func checkEncoding(encs []string, enc string) bool {
	for _, e := range encs {
		if strings.Contains(e, enc) {
			return true
		}
	}

	return false
}

func (h *handler) checkAuthCookies(cookies []*http.Cookie) (*smodels.Claims, error) {
	var errs []error
	for _, cookie := range cookies {
		claims, err := h.service.ValidateToken(cookie.Value)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		return claims, nil
	}

	return nil, fmt.Errorf("errors in token validation: %w", errors.Join(errs...))
}
