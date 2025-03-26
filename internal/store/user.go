package store

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password password `json:"-"` //changes from string to password struct
	Created  string   `json:"created_at"`
}

/*
now password will contain plain password and the hash,
hashing is done in below Hash method
*/
type password struct {
	text *string `json:"-"`
	hash []byte  `json:"-"`
}

//Hash
/*
below method gets the password string from handler class
and stores both plain pwd and hash in the password struct
*/
func (p *password) Hash(pswd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//store the password and hash in below 2 lines
	p.text = &pswd
	p.hash = hash
	return nil
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

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string) error {
	//transaction wrapper - it has 2 tasks
	//create the user
	//create the user invite
	//if one of them fails rollback both the transaction with sql transaction

	return nil
}
