package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Created  string `json:"created_at"`
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
		&user.Created)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	var userVar User
	query := `SELECT id,username,email,password,created_at FROM users WHERE id = $1`
	err := u.db.QueryRowContext(ctx, query, id).Scan(&userVar.ID, &userVar.Username, &userVar.Email, &userVar.Password, &userVar.Created)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}

	}
	return &userVar, nil
}
