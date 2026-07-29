package api

import "net/http"

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.health)

	return s.recoverMiddleware(s.requestLogger(mux))
}
