package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	AccessLevelFull  = float32(0)
	AccessLevelLeast = float32(10)

	AccessLevelSuper   = float32(1)
	AccessLevelAdmin   = float32(3)
	AccessLevelStaff   = float32(7)
	AccessLevelUser    = float32(9)
	AccessLevelVisitor = float32(9.5)
)

type Group struct {
	AccessLevel float32
}

type AuthMSClaim struct {
	UsrID string
	Group Group
	jwt.StandardClaims
}

func NewClaim(issuer, usrID string, g Group, validity time.Duration) *AuthMSClaim {
	issue := time.Now()
	expiry := issue.Add(validity)
	return &AuthMSClaim{
		UsrID: usrID,
		Group: g,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issue.Unix(),
			ExpiresAt: expiry.Unix(),
			Issuer:    issuer,
		},
	}
}
