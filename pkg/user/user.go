package user

import "time"

type User struct {
	ID          string
	Name        string
	Gender      string
	ICEPhone    string
	AvatarURL   string
	Bio         string
	Rating      float32
	NumRaters   int64
	Created     time.Time
	LastUpdated time.Time
}

type StringUpdate struct {
	IsUpdating bool
	NewValue   string
}

type UserUpdate struct {
	UserID    string
	Name      StringUpdate
	ICEPhone  StringUpdate
	Gender    StringUpdate
	AvatarURL StringUpdate
	Bio       StringUpdate
	Time      time.Time
}
