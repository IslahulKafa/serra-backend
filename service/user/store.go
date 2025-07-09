package user

import (
	"database/sql"
	"errors"
	"serra/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(u *types.User) error {
	var exists int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, u.Email).Scan(&exists)
	if err != nil {
		return err
	}

	if exists > 0 {
		return errors.New("email already registered")
	}

	err = s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE uername = ?`, u.Username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("username already taken")
	}

	res, err := s.db.Exec(`INSERT INTO users (username, email, password) VALUES (?, ?, ?)`, u.Username, u.Email, u.Password)
	if err != nil {
		return err
	}

	u.ID, _ = res.LastInsertId()
	return nil
}
