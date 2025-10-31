package handler

import (
	"net/http"
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
