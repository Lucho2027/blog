package api

import (
	"net/http"
	"time"

	"github.com/lucho2027/blog/internal/database/sqlc"
)

type PublicPost struct {
	Title     string    `json:"title"`
	Excerpt   string    `json:"excerpt"`
	Published time.Time `json:"published"`
}
type PostResponse struct {
	Posts []PublicPost `json:"posts"`
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
		Posts: make([]PublicPost, 0),
	}

	for _, p := range posts {
		resp.Posts = append(resp.Posts, PublicPost{
			Title:     p.Title,
			Excerpt:   p.Excerpt.String,
			Published: p.PublishedAt.Time,
		})
	}

	if err := writeJSON(w, http.StatusOK, resp); err != nil {
		s.log.Err(err).Msg("failed to write posts response")
	}
}
