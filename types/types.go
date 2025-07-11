package types

import "time"

type UserStore interface {
	CreateUser(u *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	UpsertPrekeyBundle(userID int64, identityKey, signedPrekey, signature string, oneTimePrekeys []string) error
	GetPrekeyBundle(userID int64) (map[string]any, error)
	SetUserProfile(userID int64, username, profilePic string) error
	SaveRefreshToken(userID int64, token string, expires time.Time) error
	GetRefreshToken(token string) (int64, error)
}

type User struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	ProfilePic string `json:"profile_pic"`
	Email      string `json:"email"`
	Password   string `json:"-"`
}
