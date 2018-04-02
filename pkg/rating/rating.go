package rating

import (
	"github.com/tomogoma/crdb"
	"github.com/tomogoma/go-typed-errors"
	"time"
)

type Rating struct {
	ID          string
	ForSection  string
	ForUserID   string
	ByUserID    string
	Rating      int32
	Comment     string
	Created     time.Time
	LastUpdated time.Time
}

type Filter struct {
	ForSection *crdb.Comparison
	ForUserID  *crdb.Comparison
	ByUserID   *crdb.Comparison
	Offset     int64
	Count      int32
}

type AverageUser struct {
	UserID    string
	Rating    float32
	NumRaters int64
}

func (f Filter) Validate() error {
	if f.ForUserID == nil && f.ByUserID == nil {
		return errors.NewClient("one of ForUserID or ByUserID must be provided")
	}
	if f.Offset < 0 {
		return errors.NewClientf("Offset must be >= 0")
	}
	if f.Count < 1 {
		return errors.NewClientf("Count must be > 0")
	}
	return nil
}
