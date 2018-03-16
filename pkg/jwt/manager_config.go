package jwt

import "github.com/tomogoma/go-typed-errors"

type ManagerConfig struct {
	jwter    JWTEr
	hs256Key []byte
}

func WithJWTEr(jwter JWTEr) Option {
	return func(conf *ManagerConfig) error {
		if jwter == nil {
			return errors.Newf("JWTEr was nil")
		}
		conf.jwter = jwter
		return nil
	}
}

func WithHS256Key(key []byte) Option {
	return func(conf *ManagerConfig) error {
		if len(key) == 0 {
			return errors.Newf("HS256 Key was empty")
		}
		conf.hs256Key = key
		return nil
	}
}
