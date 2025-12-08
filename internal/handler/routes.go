package handler

func (s *Server) addRoutes() {
	s.mux.Use(logMiddlware)
	s.mux.Use(compressMiddlware)
	s.mux.Use(extractMiddlware)

	s.mux.HandleFunc("/", s.h.errorHandler)

	s.mux.Get("/{url_id}", s.h.getFullURL)
	s.mux.Get("/ping", s.h.getPing)
	s.mux.Get("/api/user/urls", s.h.getAPIUserURLs)

	s.mux.
		With(s.h.authCookieMiddlware).
		Post("/", s.h.postURL)
	s.mux.
		With(s.h.authCookieMiddlware).
		Post("/api/shorten", s.h.postAPIShorten)
	s.mux.
		With(s.h.authCookieMiddlware).
		Post("/api/shorten/batch", s.h.postAPIShortenBatch)
}
