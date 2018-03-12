package api

import "time"

type Key struct {
	ID          string
	UserID      string
	Val         []byte
	Created     time.Time
	LastUpdated time.Time
}

func (k Key) Value() []byte {
	return k.Val
}
