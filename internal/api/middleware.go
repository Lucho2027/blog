package api

import (
	"net/http"
)

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("incoming request")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error().Any("panic", r).Msg("panic recovered")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))

			}
		}()

		next.ServeHTTP(w, r)
	})
}
