package http

import (
	"github.com/tomogoma/usersms/pkg/rating"
	"time"
)

/**
 * @apiDefine RatingsList200
 *
 * @apiSuccess (200 JSON Response) {Object[]} ratings List of ratings (values indented below).
 * @apiSuccess (200 JSON Response) {String} ratings.ID Unique identifier of this rating.
 * @apiSuccess (200 JSON Response) {String} ratings.forUserID Ratee' userID.
 * @apiSuccess (200 JSON Response) {String} ratings.byUserID Rater's userID.
 * @apiSuccess (200 JSON Response) {String} ratings.comment
 * @apiSuccess (200 JSON Response) {Integer{1-5}} ratings.rating Rating awarded by rater to ratee.
 * @apiSuccess (200 JSON Response) {String} ratings.created ISO8601 date of rating creation.
 * @apiSuccess (200 JSON Response) {String} ratings.lastUpdated Last ISO8601 date of update.
 */
type Rating struct {
	ID          string `json:"ID,omitempty"`
	ForUserID   string `json:"forUserID,omitempty"`
	ByUserID    string `json:"byUserID,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Rating      int32  `json:"rating,omitempty"`
	Created     string `json:"created,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

func NewRating(r *rating.Rating) *Rating {
	if r == nil {
		return nil
	}
	return &Rating{
		ID:          r.ID,
		ForUserID:   r.ForUserID,
		ByUserID:    r.ByUserID,
		Comment:     r.Comment,
		Rating:      r.Rating,
		Created:     r.Created.Format(time.RFC3339),
		LastUpdated: r.LastUpdated.Format(time.RFC3339),
	}
}

func NewRatings(rs []rating.Rating) []Rating {
	if len(rs) == 0 {
		return nil
	}
	var retRs []Rating
	for _, r := range rs {
		retR := NewRating(&r)
		retRs = append(retRs, *retR)
	}
	return retRs
}
