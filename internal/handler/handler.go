// Package handler implements http-server and handler for endpoints.
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/darrior/urlshortener/internal/models"
	"github.com/darrior/urlshortener/internal/service"
)

type handler struct {
	service service.IService
}

func (h *handler) errorHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Invalid request", http.StatusBadRequest)

}

func (h *handler) getFullURL(res http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("url_id")
	fullURL, err := h.service.GetURL(req.Context(), shortURL)
	if err != nil {
		http.Error(res, "Short URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}

func (h *handler) getPing(res http.ResponseWriter, req *http.Request) {
	if err := h.service.Ping(req.Context()); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *handler) postURL(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.Header.Get("content-type"), "text/plain") {
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

	shortURL, err := h.service.AddURL(req.Context(), longURL)
	if err != nil {
		http.Error(res, "Error while creating short URL", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(shortURL)))
	res.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprint(res, shortURL)
}

func (h *handler) postAPIShorten(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.Header.Get("content-type"), "application/json") {
		http.Error(res, `Content type must be "application/json"`, http.StatusBadRequest)
		return
	}

	var reqData models.ShortenerRequest

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqData); err != nil {
		http.Error(res, "Can not unmarshal JSON", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(reqData.URL); err != nil {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.AddURL(req.Context(), reqData.URL)
	if err != nil {
		http.Error(res, "Error while creating short URL", http.StatusInternalServerError)
		return
	}

	resData := models.ShortenerResponse{
		Result: shortURL,
	}

	data, err := json.Marshal(resData)
	if err != nil {
		http.Error(res, "Can not marshal short URL", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.Header().Set("content-length", strconv.Itoa(len(data)))
	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write(data)
}
