package api

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/lucho2027/blog/internal/database/sqlc"
)

type PublicPost struct {
	Title     string     `json:"title"`
	Excerpt   *string    `json:"excerpt"`
	Published *time.Time `json:"published"`
}
type PostsResponse struct {
	Posts []PublicPost `json:"posts"`
}

func parsePostsQueryParams(r *http.Request) (sqlc.GetAllPostsParams, error) {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return sqlc.GetAllPostsParams{}, err
	}
	if limit < 1 || limit > 100 {
		return sqlc.GetAllPostsParams{}, errors.New("limit needs to be a number between 1 and 100")
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return sqlc.GetAllPostsParams{}, err
	}
	if offset < 0 || offset > math.MaxInt32 {
		return sqlc.GetAllPostsParams{}, errors.New("offset is out of range")
	}
	return sqlc.GetAllPostsParams{
		OffsetVal: int32(offset),
		LimitVal:  int32(limit),
	}, nil
}

func (s *Server) posts(w http.ResponseWriter, r *http.Request) {
	p, err := parsePostsQueryParams(r)
	if err != nil {
		s.log.Err(err).Msg("failed to query posts")
		if err := writeError(w, http.StatusBadRequest, "invalid pagination parameters"); err != nil {
			s.log.Err(err).Msg("failed to write error on querying posts")
		}
		return
	}

	posts, err := s.db.Queries().GetAllPosts(r.Context(), p)
	if err != nil {
		s.log.Err(err).Msg("failed to get posts")
		if err := writeError(w, http.StatusInternalServerError, "internal server error"); err != nil {
			s.log.Err(err).Msg("failed to write error on posts")
		}
		return
	}
	resp := PostsResponse{
		Posts: make([]PublicPost, 0),
	}

	for _, post := range posts {
		var excerpt *string
		if post.Excerpt.Valid {
			excerpt = &post.Excerpt.String
		}
		var published *time.Time
		if post.PublishedAt.Valid {
			published = &post.PublishedAt.Time
		}
		resp.Posts = append(resp.Posts, PublicPost{
			Title:     post.Title,
			Excerpt:   excerpt,
			Published: published,
		})
	}

	if err := writeJSON(w, http.StatusOK, resp); err != nil {
		s.log.Err(err).Msg("failed to write posts response")
	}
}
