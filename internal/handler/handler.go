// Package handler implements http-server and handler for endpoints.
package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/darrior/urlshortener/internal/service"
)

type handler struct {
	service service.IService
}

func (h *handler) errorHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Invalid request", http.StatusBadRequest)

}

func (h *handler) postURL(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("content-type") != "text/plain" {
		http.Error(res, `Content type must be "text/plain"`, http.StatusBadRequest)
		return
	}

	rawURL, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Can not read body content", http.StatusBadRequest)
		return
	}

	longURL := string(rawURL)
	if _, err := url.ParseRequestURI(longURL); err != nil {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortURLID, err := h.service.AddURL(longURL)
	if err != nil {
		http.Error(res, "Error while creating short URL", http.StatusInternalServerError)
		return
	}

	scheme := "http://"
	if req.TLS != nil {
		scheme = "https://"
	}

	shortURL, err := url.JoinPath(scheme+req.Host, shortURLID)
	if err != nil {
		http.Error(res, "Can not join base URL with ID", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(shortURL)))
	_, _ = fmt.Fprint(res, shortURL)
}

func (h *handler) getFullURL(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("url_id")
	fullURL, err := h.service.GetURL(shortURL)
	if err != nil {
		http.Error(res, "Short URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}
