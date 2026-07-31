package api

import (
	"net/http"

	"github.com/lucho2027/blog/internal/database/sqlc"
)

type PostResponse struct {
	Posts []sqlc.Post `json:"posts"`
}

func (s *Server) posts(w http.ResponseWriter, r *http.Request) {
	p := sqlc.GetAllPostsParams{
		OffsetVal: 0,
		LimitVal:  10,
	}

	posts, err := s.db.Queries().GetAllPosts(r.Context(), p)
	if err != nil {
		s.log.Err(err).Msg("failed to get posts")
		if err := writeError(w, http.StatusInternalServerError, "Internal server error"); err != nil {
			s.log.Err(err).Msg("failed to write error on posts")
		}
		return
	}
	resp := PostResponse{
		Posts: make([]sqlc.Post, 0),
	}

	if posts == nil {
		resp.Posts = []sqlc.Post{}
	} else {
		resp.Posts = posts
	}

	if err := writeJSON(w, http.StatusOK, resp); err != nil {
		s.log.Err(err).Msg("failed to write posts response")
	}
}
