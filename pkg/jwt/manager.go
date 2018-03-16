package jwt

import (
	"github.com/tomogoma/go-typed-errors"
	"github.com/dgrijalva/jwt-go"
	jwtH "github.com/tomogoma/jwt"
)

type JWTEr interface {
	errors.IsAuthErrChecker
	Generate(claims jwt.Claims) (string, error)
	Validate(JWT string, claims jwt.Claims) (*jwt.Token, error)
}

type Manager struct {
	errors.AuthErrCheck
	JWTEr
}

type Option func(*ManagerConfig) error

func NewManager(opts ...Option) (*Manager, error) {

	conf := &ManagerConfig{}
	for i, optFunc := range opts {
		if optFunc == nil {
			return nil, errors.Newf("received nil Option at index %d", i)
		}
		if err := optFunc(conf); err != nil {
			return nil, err
		}
	}

	if conf.jwter == nil {
		var err error
		if conf.jwter, err = jwtH.NewHandler(conf.hs256Key); err != nil {
			return nil, errors.New("provide hs256Key using WithHS256Key" +
				" in order to use the default JWTEr")
		}
	}

	return &Manager{JWTEr: conf.jwter}, nil
}

func (v Manager) IsOwnerOrJWTHasAccess(JWT string, owner string, acl float32) (*Claim, error) {
	clm, err := v.JWTValid(JWT)
	if err != nil {
		return nil, err
	}
	if clm.UsrID == owner {
		return clm, nil
	}
	if err := claimsHaveAccess(*clm, acl); err != nil {
		return nil, err
	}
	return clm, nil
}

func (v Manager) JWTHasAccess(JWT string, acl float32) (*Claim, error) {
	clm, err := v.JWTValid(JWT)
	if err != nil {
		return nil, err
	}
	if err := claimsHaveAccess(*clm, acl); err != nil {
		return nil, err
	}
	return clm, nil
}

func (v Manager) JWTValid(JWT string) (*Claim, error) {
	clm := &Claim{}
	if _, err := v.JWTEr.Validate(JWT, clm); err != nil {
		if v.JWTEr.IsUnauthorizedError(err) {
			return nil, errors.NewUnauthorized(err)
		}
		if v.JWTEr.IsForbiddenError(err) {
			return nil, errors.NewForbidden(err)
		}
		if v.JWTEr.IsAuthError(err) {
			return nil, errors.NewAuth(err)
		}
		return nil, err
	}
	return clm, nil
}

func claimsHaveAccess(clms Claim, acl float32) error {
	if err := aclValid(acl); err != nil {
		return err
	}
	if clms.Group.AccessLevel > acl {
		return errors.NewForbiddenf("lack sufficient privilege to access this resource")
	}
	return nil
}

func aclValid(acl float32) error {
	if acl < AccessLevelFull || acl > AccessLevelLeast {
		return errors.Newf("invalid access level")
	}
	return nil
}
