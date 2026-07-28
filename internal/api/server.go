package api

import (
	"github.com/lucho2027/blog/internal/database"
	"github.com/rs/zerolog"
)

type Server struct {
	db  *database.DB
	log zerolog.Logger
}

func NewServer(db *database.DB, logger zerolog.Logger) *Server {
	return &Server{
		db:  db,
		log: logger,
	}
}
