package types

type UserStore interface {
	CreateUser(u *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
