// Package handler implements http-server and handler for endpoints.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/darrior/urlshortener/internal/models/api"
	"github.com/darrior/urlshortener/internal/service"
	"github.com/rs/zerolog/log"
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
	if errors.Is(err, service.ErrorURLGone) {
		res.WriteHeader(http.StatusGone)
		return
	} else if err != nil {
		http.Error(res, "Short URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, fullURL, http.StatusTemporaryRedirect)
}

func (h *handler) getPing(res http.ResponseWriter, req *http.Request) {
	if err := h.service.Ping(req.Context()); err != nil {
		log.Error().Err(err).Msg("Error while ping DB")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *handler) getAPIUserURLs(res http.ResponseWriter, req *http.Request) {
	userID := ""
	if id := req.Context().Value(_contextUserID); id != nil {
		userID = id.(string)
	}

	userURLs, err := h.service.GetUserURLs(req.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("can not get user URLs")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(userURLs) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := json.Marshal(userURLs)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add("content-type", "application/json")
	res.Header().Add("content-length", strconv.Itoa(len(data)))
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(data)
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

	status := http.StatusCreated

	userID := ""
	if id := req.Context().Value(_contextUserID); id != nil {
		userID = id.(string)
	}

	shortURL, err := h.service.AddURL(req.Context(), userID, longURL)
	if errors.Is(err, service.ErrorURLExists) {
		status = http.StatusConflict
	} else if err != nil {
		log.Error().Err(err).Msg("Can not create short URL")
		http.Error(res, "Error while creating short URL", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.Header().Set("content-length", strconv.Itoa(len(shortURL)))
	res.WriteHeader(status)
	_, _ = fmt.Fprint(res, shortURL)
}

func (h *handler) postAPIShorten(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.Header.Get("content-type"), "application/json") {
		http.Error(res, `Content type must be "application/json"`, http.StatusBadRequest)
		return
	}

	var reqData api.ShortenerRequest

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqData); err != nil {
		http.Error(res, "Can not unmarshal JSON", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(reqData.URL); err != nil {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	status := http.StatusCreated

	userID := ""
	if id := req.Context().Value(_contextUserID); id != nil {
		userID = id.(string)
	}

	shortURL, err := h.service.AddURL(req.Context(), userID, reqData.URL)
	if errors.Is(err, service.ErrorURLExists) {
		status = http.StatusConflict
	} else if err != nil {
		log.Error().Err(err).Msg("Can not create short URL")
		http.Error(res, "Error while creating short URL", http.StatusInternalServerError)
		return
	}

	resData := api.ShortenerResponse{
		Result: shortURL,
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Error().Err(err).Msg("Can not marshal short URL")
		http.Error(res, "Can not marshal short URL", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.Header().Set("content-length", strconv.Itoa(len(data)))
	res.WriteHeader(status)
	_, _ = res.Write(data)
}

func (h *handler) postAPIShortenBatch(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.Header.Get("content-type"), "application/json") {
		http.Error(res, `Content type must be "application/json"`, http.StatusBadRequest)
		return
	}

	var reqData api.ShortenerBatchRequest

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqData); err != nil {
		http.Error(res, "Can not unmarshal JSON", http.StatusBadRequest)
		return
	}

	for _, entry := range reqData {
		if _, err := url.ParseRequestURI(entry.OriginalURL); err != nil {
			http.Error(res, "Invalid URL", http.StatusBadRequest)
			return
		}
	}

	userID := ""
	if id := req.Context().Value(_contextUserID); id != nil {
		userID = id.(string)
	}

	shortURLs, err := h.service.AddURLs(req.Context(), userID, reqData)
	if err != nil {
		log.Error().Err(err).Msg("Can not create short URL")
		http.Error(res, "Error while creating short URL", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(shortURLs)
	if err != nil {
		log.Error().Err(err).Msg("Can not marshal short URL")
		http.Error(res, "Can not marshal short URL", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.Header().Set("content-length", strconv.Itoa(len(data)))
	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write(data)
}

func (h *handler) deleteAPIUserURLs(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.Header.Get("content-type"), "application/json") {
		http.Error(res, `Content type must be "application/json"`, http.StatusBadRequest)
		return
	}

	var reqData []string
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&reqData); err != nil {
		http.Error(res, "Can not unmarshal JSON", http.StatusBadRequest)
		return
	}

	log.Debug().Any("ids", reqData).Msg("IDs received")

	userID := ""
	if id := req.Context().Value(_contextUserID); id != nil {
		userID = id.(string)
	}

	if err := h.service.RemoveURLs(context.Background(), userID, reqData); err != nil {
		log.Error().Err(err).Msg("Can not delete user URLs")
		http.Error(res, "Error while deleting URLs", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}
