package http

import (
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/go-typed-errors"
)

type Config struct {
	BaseURL        string
	AllowedOrigins []string
	Guard          Guard
	Logger         logging.Logger
	Rater          Rater
	UserProfiler   UserProfiler
}

func (c Config) Validate() error {
	if c.Guard == nil {
		return errors.Newf("Guard was nil")
	}
	if c.Logger == nil {
		return errors.Newf("Logger was nil")
	}
	if c.Rater == nil {
		return errors.Newf("Rater was nil")
	}
	if c.UserProfiler == nil {
		return errors.Newf("UserProfiler was nil")
	}
	return nil
}
