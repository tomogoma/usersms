package user

import (
	"github.com/tomogoma/go-typed-errors"
	"time"
	"net/url"
	"github.com/tomogoma/usersms/pkg/jwt"
)

var validGenders = []string{"MALE", "FEMALE", "OTHER"}

type DB interface {
	errors.IsNotFoundErrChecker

	UpsertUser(UserUpdate) (*User, error)
	User(userID string, offsetUpdateDate time.Time) (*User, error)
}

type JWTEr interface {
	errors.IsAuthErrChecker
	IsOwnerOrJWTHasAccess(JWT string, owner string, acl float32) (*jwt.AuthMSClaim, error)
	JWTValid(JWT string) (*jwt.AuthMSClaim, error)
}

type FormatValidPhoner interface {
	FormatValidPhone(number string) (string, error)
}

type Manager struct {
	errors.ErrToHTTP

	db    DB
	jwter JWTEr
	pf    FormatValidPhoner
}

func NewManager(db DB, jwter JWTEr, pf FormatValidPhoner) (*Manager, error) {
	if db == nil {
		return nil, errors.Newf("nil DB")
	}
	if jwter == nil {
		return nil, errors.Newf("nil JWTEr")
	}
	if pf == nil {
		return nil, errors.Newf("nil FormatValidPhoner")
	}
	return &Manager{db: db, jwter: jwter, pf: pf}, nil
}

func (m *Manager) Update(JWT string, update UserUpdate) (*User, error) {

	_, err := m.jwter.IsOwnerOrJWTHasAccess(JWT, update.UserID, jwt.AccessLevelStaff)
	if err != nil {
		return nil, m.parseJWTErError(err, "validate JWT belongs to"+
			" subject or has access")
	}

	if err := m.validateUserUpdate(&update); err != nil {
		return nil, err
	}

	update.Time = time.Now()
	usr, err := m.db.UpsertUser(update)
	if err != nil {
		return nil, errors.Newf("upsert user: %v", err)
	}
	return usr, nil
}

func (m *Manager) User(JWT, ID string, offsetUpdateDate time.Time) (*User, error) {

	if _, err := m.jwter.JWTValid(JWT); err != nil {
		return nil, m.parseJWTErError(err, "check JWT valid")
	}

	usr, err := m.db.User(ID, offsetUpdateDate)
	if err != nil {
		if m.db.IsNotFoundError(err) {
			return nil, errors.NewNotFound("user not found")
		}
		return nil, errors.Newf("fetch user: %v", err)
	}
	return usr, nil
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

// validateUserUpdate validates r with side-effects on the ICEPhone value
// which is also formatted if valid.
func (m *Manager) validateUserUpdate(uu *UserUpdate) error {
	if uu == nil {
		return errors.Newf("validating nil user")
	}

	// UserID is required.
	if uu.UserID == "" {
		return errors.NewClient("UserID was empty")
	}

	// Name must not be empty when updating.
	if uu.Name.IsUpdating && uu.Name.NewValue == "" {
		return errors.NewClient("name was empty")
	}

	// Phone must be valid if updating and new value not empty.
	// This block also formats the newValue if valid.
	if uu.ICEPhone.IsUpdating && uu.ICEPhone.NewValue != "" {
		var err error
		uu.ICEPhone.NewValue, err = m.pf.FormatValidPhone(uu.ICEPhone.NewValue)
		if err != nil {
			return errors.NewClient(err)
		}
	}

	// Gender must be within valid values if updating.
	if uu.Gender.IsUpdating && !in(uu.Gender.NewValue, validGenders) {
		return errors.NewClientf("invalid gender value, must be one of %v",
			validGenders)
	}

	// AvatarURL must be valid if updating and not empty.
	if uu.AvatarURL.IsUpdating && uu.AvatarURL.NewValue != "" {
		if _, err := url.Parse(uu.AvatarURL.NewValue); err != nil {
			return errors.NewClientf("invalid AvatarURL: %v", err)
		}
	}

	return nil
}

func in(needle string, haystack []string) bool {
	for _, straw := range haystack {
		if needle == straw {
			return true
		}
	}
	return false
}
