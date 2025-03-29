package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// ErrNoRows
/*
this var will contain all the SQL related error
*/
var (
	ErrNoRows               = errors.New("no rows in result set")
	ErrSerializationFailure = errors.New("serialization failure")
	ErrUniqueViolation      = errors.New("unique constraint violation")
	ErrDeadlock             = errors.New("deadlock detected")
	ErrDuplicateEmail       = errors.New("a user with that email already exists")
	ErrDuplicateUsername    = errors.New("a user with that username already exists")

	//will timeout the query execution post 10 sec, used in all store classes
	QueryTimeout = time.Second * 10
)

/*
Storage is a wrapper struct that groups different
store components (Post & User) mostly based on DB tables.
It does not implement any database logic itself;
instead, it holds interfaces for
PostStore and UserStore Operations like create, insert, update based on use case.
*/
type Storage struct {
	Post interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, post *Post) (error, *Post)
		//pagination struct added here as change
		GetUserFeed(ctx context.Context, id int64, pagination Pagination) ([]*PostWithMetaData, error)
	}
	User interface {
		Create(ctx context.Context, tx *sql.Tx, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
		//created for invite a user and create it
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(ctx context.Context, token string) error
	}

	Comment interface {
		Create(ctx context.Context, comment *Comment) error
		GetCommentOfUserOnPost(ctx context.Context, postId int64) ([]*Comment, error)
	}

	Follow interface {
		Follow(ctx context.Context, userId, followingId int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	/*	this method to be called in api class and main class
		and value to be set in api.go, from there it will go to handler and inside handler
		these to be called.

		As create method of both PostStore and UserStore
		use *pointer for Create and to satisfy interface concept for Post
		and User Interface, we need to pass &
	*/
	return Storage{
		Post:    &PostStore{db: db},
		User:    &UserStore{db: db},
		Comment: &CommentStore{db: db},
		Follow:  &FollowerStore{db: db},
	}
}

/*
Transaction wrapper function which will be used in user store
It's used to roll back the transaction if an error occurs
*/
func withTxn(db *sql.DB, ctx context.Context, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//now we will use this txn in the fn
	if err := f(tx); err != nil {
		_ = tx.Rollback() //ignoring the error as of now
		return err
	}
	err = tx.Commit() //please ensure to commit or it wont reflect in DB
	if err != nil {
		return err
	}

	return nil
}
