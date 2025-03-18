package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Follower struct {
	UserId     int64  `json:"user_id"`
	FollowerId int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

//Follow
/*
even though the follower handler is written in user class,
for follow and unfollow we decided to go with FollowerStore class
separately
*/
func (f *FollowerStore) Follow(ctx context.Context, userId, followingId int64) error {
	query := `INSERT INTO followers(user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, userId, followingId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNoRows
		case err.(*pq.Error) != nil:
			var pqErr *pq.Error
			errors.As(err, &pqErr) //read about errors.As in GPT
			switch pqErr.Code {
			case "40001": // Serialization failure (retry recommended)
				return ErrSerializationFailure
			case "23505": // Unique constraint violation
				return ErrUniqueViolation
			case "40P01": // Deadlock detected
				return ErrDeadlock
			default:
				return err
			}
		default:
			return err
		}

	}
	return nil
}
