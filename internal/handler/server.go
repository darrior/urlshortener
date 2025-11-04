package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/darrior/urlshortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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

func (s *Server) Run(ctx context.Context) error {

	go func() {
		<-ctx.Done()
		if err := s.srv.Close(); err != nil {
			log.Error().Err(err).Msg("Can not stop server properly")
		}
		log.Info().Msg("Server shutdown gracefuly")
	}()

	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
