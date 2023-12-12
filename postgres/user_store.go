package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/noloman/goreddit"
)

// UserStore is a struct with a reference to sqlx.DB
type UserStore struct {
	DB *sqlx.DB
}

func (s *UserStore) User(id uuid.UUID) (goreddit.User, error) {
	var u goreddit.User
	if err := s.DB.Get(&u, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return goreddit.User{}, fmt.Errorf("Error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) UserByUsername(username string) (goreddit.User, error) {
	var u goreddit.User
	if err := s.DB.Get(&u, `SELECT * FROM users WHERE username = $1`, username); err != nil {
		return goreddit.User{}, fmt.Errorf("Error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) Users() ([]goreddit.User, error) {
	var uu []goreddit.User
	if err := s.DB.Select(&uu, `SELECT * FROM users`); err != nil {
		return []goreddit.User{}, fmt.Errorf("Error getting users: %w", err)
	}
	return uu, nil
}

func (s *UserStore) CreateUser(t *goreddit.User) error {
	if err := s.DB.Get(t, `INSERT INTO users VALUES ($1, $2, $3) RETURNING *`,
		t.ID,
		t.Username,
		t.Password); err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}
	return nil
}

func (s *UserStore) UpdateUser(t *goreddit.User) error {
	if err := s.DB.Get(t, `UPDATE users SET username = $1, password = $2 WHERE id = $3 RETURNING`,
		t.Username,
		t.Password,
		t.ID); err != nil {
		return fmt.Errorf("Error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.DB.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("Error deleting user: %w", err)
	}
	return nil
}
