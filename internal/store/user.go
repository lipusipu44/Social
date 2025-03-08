package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Created  int64  `json:"created_at"`
}

//UserStore
/*
This struct contains the DB connection which is to be used in below methods
as UserStore is passed as reference in the methods
*/
type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
				INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
				RETURNING id,created_at
				`
	err := s.db.QueryRowContext(ctx,
		query,
		user.Email,
		user.Username,
		user.Password,
	).Scan(&user.ID,
		&user.Created,
		&user.Password)
	if err != nil {
		return err
	}
	return nil
}
