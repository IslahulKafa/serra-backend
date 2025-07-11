package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"serra/types"
	"time"
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
	err := s.db.QueryRow(`SELECT id, username, email, password FROM users WHERE email = ?`, email).Scan(&u.ID, &u.Username, &u.ProfilePic, &u.Email, &u.Password)
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
		Scan(&u.ID, &u.Username, &u.ProfilePic, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Store) UpsertPrekeyBundle(userID int64, identityKey, signedPrekey, signature string, oneTimePrekeys []string) error {
	prekeysJSON, err := json.Marshal(oneTimePrekeys)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO prekeys (user_id, identity_key, signed_prekey, signed_prekey_signature, one_time_prekeys)
	VALUES (?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	identity_key = VALUES(identity_key),
	signed_prekey = VALUES(signed_prekey),
	signed_prekey_signature = VALUES(signed_prekey_signature),
	one_time_prekeys = VALUES(one_time_prekeys)`, userID, identityKey, signedPrekey, signature, prekeysJSON)

	return err
}

func (s *Store) GetPrekeyBundle(userID int64) (map[string]any, error) {
	var (
		identityKey  string
		signedPrekey string
		signature    string
		prekeysJSON  string
	)

	err := s.db.QueryRow(`SELECT identity_key, signed_prekey, signed_prekey_signature, one_time_prekeys
	FROM prekeys
	WHERE user_id = ?`, userID).Scan(&identityKey, &signedPrekey, &signature, &prekeysJSON)
	if err != nil {
		return nil, err
	}

	// Decode json array
	var prekeys []string
	if err := json.Unmarshal([]byte(prekeysJSON), &prekeys); err != nil {
		return nil, err
	}

	if len(prekeys) == 0 {
		return nil, errors.New("no one-time prekeys available")
	}

	// Pop one prekeys
	oneTimePrekey := prekeys[0]
	remaining := prekeys[1:]

	// Update, remaining list in db
	updated, err := json.Marshal(remaining)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`UPDATE prekeys SET one_time_prekeys = ? WHERE user_id = ?`, updated, userID)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"identity_key":            identityKey,
		"signed_prekey":           signedPrekey,
		"signed_prekey_signature": signature,
		"ne_time_prekey":          oneTimePrekey,
	}, nil
}

func (s *Store) SetUserProfile(userID int64, username, profilePic string) error {
	var exists int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ? AND id != ?`, username, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("Username already taken")
	}

	_, err = s.db.Exec(`UPDATE users SET username = ?, profile_pic = ? WHERE id = ?`, username, profilePic, userID)

	return err
}

func (s *Store) SaveRefreshToken(userID int64, token string, expires time.Time) error {
	_, err := s.db.Exec(`INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES (?, ?, ?)`, token, userID, expires)
	return err
}

func (s *Store) GetRefreshToken(token string) (int64, error) {
	var userID int64
	var expiresAt time.Time
	err := s.db.QueryRow(`SELECT user_id, expires_at FROM refresh_tokens WHERE token = ?`, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("refresh token expired")
	}

	return userID, nil
}
