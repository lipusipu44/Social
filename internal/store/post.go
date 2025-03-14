package store

import (
	"context"
	"database/sql"
	"errors"
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
	ID      int64      `json:"id"`
	Content string     `json:"content"`
	Title   string     `json:"title"`
	UserID  int64      `json:"user_id"`
	Tags    []string   `json:"tags"`
	Created string     `json:"created_at"`
	Updated string     `json:"updated_at"`
	Comment []*Comment `json:"comments"` //comment not part of post table, just for showing comments on posts
}

func (p *PostStore) Create(ctx context.Context, post *Post) error {
	/*
		ctx : It's a context object used to handle timeouts, cancellations,
		and request-scoping.
		It ensures that the database operation doesn’t run indefinitely.

		these $1,$2... to be fetched from post struct which is to be passed
		from posts.go in cmd/api package
		from handler method
	*/
	query := `INSERT INTO posts (content, title, user_id,tags)
values ($1, $2, $3, $4) RETURNING id,created_at,updated_at`
	/*
		above return part is used in scan section below
	*/
	/*
		$1,$2... order is called in the same way mentioned in INSERT query above
		in below line.
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

func (p *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id,title,user_id,content,created_at,tags,updated_at FROM posts WHERE id = $1`
	var postVar Post
	err := p.db.QueryRowContext(ctx,
		query, id).Scan(
		&postVar.ID,
		&postVar.Title,
		&postVar.UserID,
		&postVar.Content,
		&postVar.Created,
		pq.Array(&postVar.Tags),
		&postVar.Updated)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, err
		}
	}
	return &postVar, nil
}
