package handler

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/darrior/urlshortener/internal/service/models"
	"github.com/rs/zerolog/log"
)

const _authCookieName = "auth_cookie"

type contextUserID string

const _contextUserID contextUserID = "user_id"

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

func (h *handler) authCookieMiddlware(n http.Handler) http.Handler {
	authCookieHandler := func(res http.ResponseWriter, req *http.Request) {
		cookies := req.CookiesNamed(_authCookieName)

		claims, err := h.checkAuthCookies(cookies)
		if err != nil || claims.UserID == "" {
			userID, err := h.service.NewUserID()
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			claims = &models.Claims{
				UserID: userID,
			}
		}

		userID := claims.UserID

		tokenString, err := h.service.SignClaims(claims)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		nextReq := req.WithContext(context.WithValue(req.Context(), _contextUserID, userID))

		n.ServeHTTP(res, nextReq)

		cookie := &http.Cookie{
			Name:  _authCookieName,
			Value: tokenString,
		}

		http.SetCookie(res, cookie)
	}

	return http.HandlerFunc(authCookieHandler)
}

func compressMiddlware(h http.Handler) http.Handler {
	compressHandler := func(res http.ResponseWriter, req *http.Request) {
		if !checkEncoding(req.Header.Values("accept-encoding"), "gzip") {
			h.ServeHTTP(res, req)
			return
		}

		var buf bytes.Buffer
		newRes := gzipResponseWriter{
			ResponseWriter: res,
			status:         0,
			writer:         &buf,
		}

		h.ServeHTTP(&newRes, req)

		if !strings.HasPrefix(res.Header().Get("content-type"), "application/json") &&
			!strings.HasPrefix(res.Header().Get("content-type"), "text/html") {
			res.WriteHeader(newRes.status)
			_, _ = res.Write(buf.Bytes())

			return
		}

		var gbuf bytes.Buffer
		w, err := gzip.NewWriterLevel(&gbuf, gzip.BestSpeed)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(buf.Bytes()); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := w.Close(); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		body := gbuf.Bytes()
		res.Header().Set("content-encoding", "gzip")
		res.Header().Del("content-length")
		res.Header().Set("content-length", strconv.Itoa(len(body)))

		res.WriteHeader(newRes.status)

		_, _ = res.Write(body)
	}

	return http.HandlerFunc(compressHandler)
}
