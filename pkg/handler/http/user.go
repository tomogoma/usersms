package http

import (
	"github.com/tomogoma/usersms/pkg/user"
)

/**
 * @apiDefine User200
 *
 * @apiSuccess (200 JSON Response) {String} ID User's ID.
 * @apiSuccess (200 JSON Response) {String} name
 * @apiSuccess (200 JSON Response) {String} ICEPhone User's (In Case of Emergency) phone number.
 * @apiSuccess (200 JSON Response) {String="MALE","FEMALE","OTHER"} gender
 * @apiSuccess (200 JSON Response) {String} avatarURL User's profile picture URL.
 * @apiSuccess (200 JSON Response) {String} bio Brief description of user.
 * @apiSuccess (200 JSON Response) {Float{0-5}} rating Overall rating of user.
 * @apiSuccess (200 JSON Response) {String} created ISO8601 date of user profile creation.
 * @apiSuccess (200 JSON Response) {String} lastUpdated last ISO8601 date when this profile was updated.
 */
type User struct {
	ID          string  `json:"ID,omitempty"`
	Name        string  `json:"name,omitempty"`
	ICEPhone    string  `json:"ICEPhone,omitempty"`
	Gender      string  `json:"gender,omitempty"`
	AvatarURL   string  `json:"avatarURL,omitempty"`
	Bio         string  `json:"bio,omitempty"`
	Rating      float32 `json:"rating,omitempty"`
	Created     string  `json:"created,omitempty"`
	LastUpdated string  `json:"lastUpdated,omitempty"`
}

func NewUser(u *user.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:        u.ID,
		Name:      u.Name,
		ICEPhone:  u.ICEPhone,
		Gender:    u.Gender,
		AvatarURL: u.AvatarURL,
		Bio:       u.Bio,
		Rating:    u.Rating,
	}
}
