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

	res, err := s.db.Exec(`INSERT INTO users (username, email, password) VALUES (?, ?, ?)`, u.Username, u.Email, u.Password)
	if err != nil {
		return err
	}

	u.ID, _ = res.LastInsertId()
	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow(`SELECT id, username, email, password FROM users WHERE email = ?`, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (s *Store) GetUserByID(id int64) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow(`SELECT id, username, email, password FROM users WHERE id = ?`, id).
		Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
