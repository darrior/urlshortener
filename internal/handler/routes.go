package handler

func (s *Server) addRoutes() {
	s.mux.Use(logMiddlware)
	s.mux.Use(extractMiddlware)

	s.mux.HandleFunc("/", s.h.errorHandler)

	s.mux.Get("/{url_id}", s.h.getFullURL)

	s.mux.With(compressMiddlware).Post("/", s.h.postURL)
	s.mux.Post("/api/shorten", s.h.postAPIShorten)
}
