package store

import (
	"context"
	"database/sql"
)

/*
Storage is a wrapper struct that groups different
store components (Post & User).
It does not implement any database logic itself;
instead, it holds interfaces for PostStorage and UserStorage.
*/
type Storage struct {
	Post interface {
		Create(ctx context.Context) error
	}
	User interface {
		Create(ctx context.Context) error
	}
}

func NewStorage(db *sql.DB) Storage {
	/*	this method to be called in api class and main class
		and value to be set in api.go, from there it will go to handler and inside handler
		these to be called.

		As create method of both PostStorage and UserStorage
		use *pointer for Create and to satisfy interface concept for Post
		and User Interface we need to pass &
	*/
	return Storage{
		Post: &PostStorage{db: db},
		User: &UserStorage{db: db},
	}
}
