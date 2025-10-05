package handler

func (s *Server) addRoutes() {
	s.mux.HandleFunc("/", s.h.errorHandler)

	s.mux.HandleFunc("GET /{url_id}", s.h.getFullURL)

	s.mux.HandleFunc("POST /", s.h.postURL)
}
