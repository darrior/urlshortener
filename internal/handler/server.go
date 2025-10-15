package handler

import (
	"net/http"

	"github.com/darrior/urlshortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	mux     *chi.Mux
	h       *handler
	address string
}

func NewServer(address string, service *service.Service) *Server {
	s := Server{
		mux: chi.NewRouter(),
		h: &handler{
			service: service,
		},
		address: address,
	}

	s.addRoutes()

	return &s
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.address, s.mux)
}
