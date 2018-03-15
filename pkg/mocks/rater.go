package mocks

import (
	"github.com/tomogoma/usersms/pkg/rating"
	"github.com/tomogoma/go-typed-errors"
)

type Rater struct {
	errors.ErrToHTTP

	RtUsrRecTkn  string
	RtUsrRecRtng rating.Rating
	RtUsrErr     error

	RtngsRecTkn  string
	RtngsRecFltr rating.Filter
	RtngsRtng    []rating.Rating
	RtngsErr     error
}

func (r *Rater) RateUser(token string, rating rating.Rating) error {
	r.RtUsrRecTkn = token
	r.RtUsrRecRtng = rating
	return r.RtUsrErr
}

func (r *Rater) Ratings(token string, filter rating.Filter) ([]rating.Rating, error) {
	r.RtngsRecTkn = token
	r.RtngsRecFltr = filter
	return r.RtngsRtng, r.RtngsErr
}
