package store

import (
	"context"
	"database/sql"
)

/*
This struct contains the DB connection which is to be used in below methods
as PostStorage is passed as reference in the methods
*/
type PostStorage struct {
	db *sql.DB
}

func (p *PostStorage) Create(ctx context.Context) error {
	return nil
}
