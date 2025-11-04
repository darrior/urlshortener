package handler

func (s *Server) addRoutes() {
	s.mux.Use(logMiddlware)
	s.mux.Use(extractMiddlware)

	s.mux.HandleFunc("/", s.h.errorHandler)

	s.mux.With(compressMiddlware).Get("/{url_id}", s.h.getFullURL)

	s.mux.Post("/", s.h.postURL)
	s.mux.With(compressMiddlware).Post("/api/shorten", s.h.postAPIShorten)
}
