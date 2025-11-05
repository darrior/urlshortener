package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/darrior/urlshortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	mux *chi.Mux
	h   *handler
	srv *http.Server
}

func NewServer(address string, service *service.Service) *Server {
	s := Server{
		mux: chi.NewRouter(),
		h: &handler{
			service: service,
		},
	}

	srv := http.Server{
		Addr:    address,
		Handler: s.mux,
	}
	s.srv = &srv
	s.addRoutes()

	return &s
}

func (s *Server) Run() error {
	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	<-ctx.Done()
	if err := s.srv.Close(); err != nil {
		return err
	}
	return nil
}
