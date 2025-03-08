package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
)

//PostStore
/*
This struct contains the DB connection which is to be used in below methods
as PostStore is passed as reference in the methods
*/
type PostStore struct {
	db *sql.DB
}

//Post
/*
📌 This struct represents a post in the application.
📌 It uses struct tags (json:"field_name") to map Go fields to JSON keys.
EX: {"id":1,"title":"Hello, Go!","content":"This is a sample post","user_id":42}
in POST Payload
📌 Tags is a slice of strings ([]string), useful for categorizing posts.
*/
type Post struct {
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	Title   string   `json:"title"`
	UserID  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
	Created int64    `json:"created_at"`
	Updated int64    `json:"updated_at"`
}

func (p *PostStore) Create(ctx context.Context, post *Post) error {
	/*
		ctx : It's a context object used to handle timeouts, cancellations,
		and request-scoping.
		It ensures that the database operation doesn’t run indefinitely.
	*/
	query := `INSERT INTO posts (content, title, user_id,tags)
values ($1, $2, $3, $4) RETURNING id,created_at,updated_at`
	/*
		above return part is used in scan section below
	*/

	err := p.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.Created,
		&post.Updated)
	/*
		Scan reads the returned values (id, created_at, updated_at) from the database.
		Stores them back into the post object.
	*/
	if err != nil {
		return err
	}
	return nil
}
