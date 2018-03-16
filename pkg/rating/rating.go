package rating

import (
	"time"
	"github.com/tomogoma/go-typed-errors"
)

type Rating struct {
	ID          string
	ForUserID   string
	ByUserID    string
	Comment     string
	Rating      int32
	Created     time.Time
	LastUpdated time.Time
}

type Filter struct {
	ForUserID string
	ByUserID  string
	Offset    int64
	Count     int32
}

func (f Filter) Validate() error {
	if f.ForUserID == "" && f.ByUserID == "" {
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
