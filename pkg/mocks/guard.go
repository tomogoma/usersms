package mocks

import (
	"github.com/tomogoma/usersms/pkg/api"
	"github.com/tomogoma/go-typed-errors"
)

type Guard struct {
	errors.AuthErrCheck

	ExpAPIKValidUsrID string
	ExpAPIKValidErr   error
	ExpNewAPIK        *api.Key
	ExpNewAPIKErr     error
}

func (g *Guard) APIKeyValid(key []byte) (string, error) {
	return g.ExpAPIKValidUsrID, g.ExpAPIKValidErr
}
func (g *Guard) NewAPIKey(userID string) (*api.Key, error) {
	return g.ExpNewAPIK, g.ExpNewAPIKErr
}
