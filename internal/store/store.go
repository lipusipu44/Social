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
	}
	User interface {
		Create(ctx context.Context, user *User) error
	}

	Comment interface {
		GetCommentOfUserOnPost(ctx context.Context, postId int64) ([]*Comment, error)
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
	}
}
