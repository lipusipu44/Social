package store

import (
	"context"
	"database/sql"
)

/*
This struct contains the DB connection which is to be used in below methods
as UserStorage is passed as reference in the methods
*/
type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context) error {
	return nil
}
