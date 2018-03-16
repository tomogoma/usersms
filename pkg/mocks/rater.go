package mocks

import (
	"github.com/tomogoma/usersms/pkg/rating"
	"github.com/tomogoma/go-typed-errors"
)

type Rater struct {
	errors.ErrToHTTP

	RtUsrRecTkn     string
	RtUsrRecFrUsrID string
	RtUsrRecCmnt    string
	RtUsrRecRtng    int32
	RtUsrErr        error

	RtngsRecTkn  string
	RtngsRecFltr rating.Filter
	RtngsRtng    []rating.Rating
	RtngsErr     error
}

func (r *Rater) RateUser(token string, forUserID, comment string, rating int32) error {
	r.RtUsrRecTkn = token
	r.RtUsrRecFrUsrID = forUserID
	r.RtUsrRecCmnt = comment
	r.RtUsrRecRtng = rating
	return r.RtUsrErr
}

func (r *Rater) Ratings(token string, filter rating.Filter) ([]rating.Rating, error) {
	r.RtngsRecTkn = token
	r.RtngsRecFltr = filter
	return r.RtngsRtng, r.RtngsErr
}
