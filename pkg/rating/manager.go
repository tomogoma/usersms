package rating

import (
	"github.com/tomogoma/go-typed-errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	jwtH "github.com/tomogoma/usersms/pkg/jwt"
)

type JWTEr interface {
	errors.IsAuthErrChecker
	JWTValidOnClaim(JWT string, clm jwt.Claims) error
	JWTValid(JWT string) (*jwtH.AuthMSClaim, error)
}

type IDEr interface {
	NextID() (string, error)
}

type DB interface {
	errors.IsNotFoundErrChecker
	SaveRating(rating Rating) error
	Rating(byUserID, forSection, forUserID string) (*Rating, error)
	Ratings(Filter) ([]Rating, error)
}

type Manager struct {
	errors.ErrToHTTP

	jwter JWTEr
	db    DB
	idgen IDEr
}

func NewManager(jwter JWTEr, db DB, idGen IDEr) (*Manager, error) {
	if jwter == nil {
		return nil, errors.Newf("nil JWTEr")
	}
	if db == nil {
		return nil, errors.Newf("nil DB")
	}
	if idGen == nil {
		return nil, errors.Newf("nil IDEr")
	}
	return &Manager{jwter: jwter, db: db, idgen: idGen}, nil
}

func (m *Manager) RateUser(JWT, forUserID, comment string, rating int32) error {

	clm, err := m.jwtCanRate(JWT)
	if err != nil {
		return err
	}

	_, err = m.db.Rating(clm.ByUsrID, clm.ForSection, forUserID)
	if err == nil {
		return errors.NewClientf("user already rated by JWT owner in JWT provided section")
	}
	if !m.db.IsNotFoundError(err) {
		return errors.Newf("")
	}

	ID, err := m.idgen.NextID()
	if err != nil {
		return errors.Newf("generate ID: %v", err)
	}

	now := time.Now()
	err = m.db.SaveRating(Rating{ID: ID, ForSection: clm.ForSection,
		ForUserID: forUserID, ByUserID: clm.ByUsrID, Rating: rating,
		Comment: comment, Created: now, LastUpdated: now})
	if err != nil {
		return errors.Newf("save rating: %v", err)
	}

	return nil
}

func (m *Manager) Ratings(JWT string, filter Filter) ([]Rating, error) {

	if _, err := m.jwter.JWTValid(JWT); err != nil {
		return nil, m.parseJWTErError(err, "check JWT valid")
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	rtngs, err := m.db.Ratings(filter)
	if err != nil {
		if m.db.IsNotFoundError(err) {
			return nil, errors.NewNotFound("ratings not found for filter")
		}
		return nil, errors.Newf("fetch ratings: %v", err)
	}

	return rtngs, nil
}

func (m *Manager) jwtCanRate(JWT string) (*Claim, error) {
	clm := &Claim{}
	if err := m.jwter.JWTValidOnClaim(JWT, clm); err != nil {
		return nil, m.parseJWTErError(err, "validate JWT on claim")
	}
	return clm, nil
}

func (m *Manager) parseJWTErError(err error, errCtx string) error {
	if m.jwter.IsAuthError(err) || m.jwter.IsUnauthorizedError(err) {
		return errors.NewUnauthorized(err)
	}
	if m.jwter.IsForbiddenError(err) {
		return errors.NewForbidden(err)
	}
	return errors.Newf("%s: %v", errCtx, err)
}
