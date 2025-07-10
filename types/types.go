package types

type UserStore interface {
	CreateUser(u *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	UpsertPrekeyBundle(userID int64, identityKey, signedPrekey, signature string, oneTimePrekeys []string) error
	GetPrekeyBundle(userID int64) (map[string]any, error)
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
