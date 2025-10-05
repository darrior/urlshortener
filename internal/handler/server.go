package handler

import (
	"net/http"

	"github.com/darrior/urlshortener/internal/service"
)

type Server struct {
	mux     *http.ServeMux
	h       *handler
	address string
}

func NewServer(address string, service *service.Service) *Server {
	s := Server{
		mux: http.NewServeMux(),
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
