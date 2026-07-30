package api

import "net/http"

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /posts", s.posts)
	return s.recoverMiddleware(s.requestLogger(mux))
}
