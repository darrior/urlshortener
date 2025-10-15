package handler

func (s *Server) addRoutes() {
	s.mux.HandleFunc("/", s.h.errorHandler)

	s.mux.Get("/{url_id}", s.h.getFullURL)

	s.mux.Post("/", s.h.postURL)
}
