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
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	Title   string   `json:"title"`
	UserID  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
	Created string   `json:"created_at"`
	Updated string   `json:"updated_at"`
	Version int64    `json:"version"` /*
		version will check the version passed in request, if matches, then
		update will happen, and it will increase the version in each update"
	*/
	Comment []*Comment `json:"comments"` //comment not part of post table, just for showing comments on posts
	User    User       `json:"user"`
}

type PostWithMetaData struct {
	Post
	Commentcount int `json:"comment_count"`
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
		here ctx is not used directly in timeout, its
		used inside and that ctx will be carry forward,

		check details of WithTimeout cancel is a func type,
		dats why we can use defer cancel()
		def:

		func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
		type CancelFunc func()
	*/
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
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
	query := `SELECT id,title,user_id,content,created_at,tags,updated_at,version FROM posts WHERE id = $1`
	/*
		here ctx is not used directly in timeout, its
		used inside and that ctx will be carry forward,

		check details of WithTimeout cancel is a func type,
		dats why we can use defer cancel()
		def:

		func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
		type CancelFunc func()
	*/
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	var postVar Post
	err := p.db.QueryRowContext(ctx,
		query, id).Scan(
		&postVar.ID,
		&postVar.Title,
		&postVar.UserID,
		&postVar.Content,
		&postVar.Created,
		pq.Array(&postVar.Tags),
		&postVar.Updated,
		&postVar.Version)
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

//Delete
/*
delete the post based on post id
*/
func (p *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	/*
		here ctx is not used directly in timeout, its
		used inside and that ctx will be carry forward,

		check details of WithTimeout cancel is a func type,
		dats why we can use defer cancel()
		def:

		func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
		type CancelFunc func()
	*/
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	res, err := p.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	//self-explanatory
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoRows
	}
	return nil
}

//Update
/*
for concurrency control, it checks the version, say 2 concurrent user
comes for update, before update in handler using middleware they would get
the Post struct, say both the user has version-1 in post struct, 1st user updated
post, version increased to 1 from 0, but second user has version-0 in his struct as
both are concurrent user who are doing update at a time, there it will fail, as db has version-1
2nd user has version-2, there it will block operation of update for second user and will ask for
retry
*/
func (p *PostStore) Update(ctx context.Context, post *Post) (error, *Post) {
	query := `
				UPDATE posts
				SET content = $1, 
				title = $2 ,version= version+1
				where id=$3 and version = $4
				RETURNING  title,content,version`

	/*
		here ctx is not used directly in timeout, its
		used inside and that ctx will be carry forward,

		check details of WithTimeout cancel is a func type,
		dats why we can use defer cancel()
		def:

		func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
		type CancelFunc func()
	*/
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	//this part only I was confused, while doing
	var postVar *Post = post

	err := p.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.ID, post.Version).Scan(&postVar.Title, &postVar.Content, &postVar.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNoRows, nil
		//copy paste from GPT, dont panic for rest of the switch-case
		case err.(*pq.Error) != nil:
			var pqErr *pq.Error
			errors.As(err, &pqErr) //read about errors.As in GPT
			switch pqErr.Code {
			case "40001": // Serialization failure (retry recommended)
				return ErrSerializationFailure, nil
			case "23505": // Unique constraint violation
				return ErrUniqueViolation, nil
			case "40P01": // Deadlock detected
				return ErrDeadlock, nil
			default:
				return err, nil
			}

		default:
			return err, nil
		}

	}
	return nil, postVar
}

//GetUserFeed
/*
here big query is used, but imp part to check is
PostMetaData struct post struct, but insted of doing
postmetadata.post.id we can do directly postmetadata.id which
will point to post's id field, this I was not aware
*/
func (p *PostStore) GetUserFeed(ctx context.Context, id int64) ([]*PostWithMetaData, error) {
	query := `
SELECT
    p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
    u.username,
    COUNT(c.id) AS comments_count
FROM posts p
LEFT JOIN comments c ON c.post_id = p.id
LEFT JOIN users u ON p.user_id = u.id
JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
WHERE f.user_id = $1 OR p.user_id = $1
GROUP BY p.id, u.username
ORDER BY p.created_at DESC;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//here bit of changes to capture the array of feeds
	var postMetaVars []*PostWithMetaData
	for rows.Next() {
		var p PostWithMetaData
		err = rows.Scan(&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.Created,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.Commentcount)
		if err != nil {
			return nil, err
		}
		postMetaVars = append(postMetaVars, &p)
	}
	return postMetaVars, nil
}
