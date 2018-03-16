package rating

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

const claimTokenValidity = 24 * 7 * time.Hour

type Claim struct {
	ByUsrID    string
	ForSection string
	jwt.StandardClaims
}

func NewClaim(issuer, byUsrID, forSection string) *Claim {
	issue := time.Now()
	expiry := issue.Add(claimTokenValidity)
	return &Claim{
		ByUsrID:    byUsrID,
		ForSection: forSection,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issue.Unix(),
			ExpiresAt: expiry.Unix(),
			Issuer:    issuer,
		},
	}
}
