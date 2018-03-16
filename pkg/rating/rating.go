package rating

import "time"

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
