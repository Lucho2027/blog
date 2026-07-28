package api

import "net/http"

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	err := s.db.Ping(r.Context())
	if err != nil {
		s.log.Err(err).Msg("database health check failed")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("service unavailable"))
		return
	}

	s.log.Debug().Msg("Health check passed")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
