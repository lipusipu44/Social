package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password password `json:"-"` //changes from string to password struct
	Created  string   `json:"created_at"`
	IsActive bool     `json:"is_active"`
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

// Create
// replaced s.db with tx in QueryRowContext
func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
				INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
				RETURNING id,created_at
				`
	err := tx.QueryRowContext(ctx,
		query,
		user.Email,
		user.Username,
		user.Password.hash, //earlier used to be plain text now its hashed version
	).Scan(&user.ID,
		&user.Created)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (u *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	var userVar User
	query := `SELECT id,username,email,password,created_at FROM users WHERE id = $1 and is_active = true`
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

//GetByEmail
/*
Simple method to use txn to get the user by checking the email
in the table
*/
func (u *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var userVar User
	query := `SELECT id,username,email,created_at
			  FROM users WHERE email = $1
			  and is_active = true`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	err = tx.QueryRowContext(ctx, query, email).Scan(&userVar.ID, &userVar.Username, &userVar.Email, &userVar.Created)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err := tx.Rollback()
			if err != nil {
				return nil, err
			}
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return &userVar, nil
}

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	//transaction wrapper - it has 2 tasks
	/*
		This below method looks confusing, but its a normal method call
		like others withTxn returns an error so we called it as per
		this func signature, in withTxn method it needs a func, with param as
		txn and that's what we did the logic for and if there is an error then txn
		roll back happens, so basically this func calls user Create and sendInvite
		both the logic and if any error then it sends the error to withTxn in Store
		.go file

		Basically, withTxn needs a function and if that function generates the error
		then it rolls back the txn accordingly
	*/
	return withTxn(u.db, ctx, func(tx *sql.Tx) error {
		//create the user
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}
		//create the user invite
		//if one of them fails rollback both the transaction with sql transaction
		if err := u.CreateUserInvitation(ctx, tx, user.ID, token, invitationExp); err != nil {
			return err
		}
		return nil
	})

}

func (u *UserStore) CreateUserInvitation(ctx context.Context, tx *sql.Tx, userID int64, token string, invitationExp time.Duration) error {
	query := `INSERT INTO user_invitations (token,user_id, expiry) VALUES ($3, $1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID, time.Now().Add(invitationExp), token)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) Activate(ctx context.Context, token string) error {
	return withTxn(u.db, ctx, func(tx *sql.Tx) error {
		//find the user that this token belongs to
		user, err := u.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		//if found update the user
		if err := u.updateUser(ctx, tx, user); err != nil {
			return err
		}
		//delete the invitation from user_invitation table
		if err := u.deleteUSerInvitations(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (u *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
				`

	/*
		below 2 lines will convert normal token to hashToken,
		as hashtoken is only stored in DB, not token string
	*/

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var userVar User
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&userVar.ID,
		&userVar.Username,
		&userVar.Email,
		&userVar.Created,
		&userVar.IsActive,
	)
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

func (u *UserStore) updateUser(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `Update users set is_active=true where id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) deleteUSerInvitations(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
