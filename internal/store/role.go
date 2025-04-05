package store

import (
	"context"
	"database/sql"
	"errors"
)

type Role struct {
	Id    int64
	Name  string
	Level string
	Desc  string
}

type RoleStore struct {
	db *sql.DB
}

//GetByName
/*
To get the role details by name, very simple way
*/
func (r *RoleStore) GetByName(name string, ctx context.Context) (*Role, error) {
	tx, _ := r.db.BeginTx(ctx, nil)
	role := &Role{}
	query := "SELECT id, name,level FROM roles WHERE name = $1"
	err := tx.QueryRowContext(ctx, query, name).Scan(role.Id, &role.Name, &role.Level)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}

	}
	return role, nil
}
