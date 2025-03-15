package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	//User is not part of comment table, its for join and data we are using, demo purpose
	User User `json:"user"` //user is fit into show username in Comment and in join also User table used
}

type CommentStore struct {
	db *sql.DB
}

//GetCommentOfUserOnPost
/*
Get the comment based on post id,
join with user to show just how join works with user

Ex:
"comments": [
        {
            "id": 2,
            "post_id": 3,
            "user_id": 1,
            "content": "I completely agree with this point.",
            "created_at": "2025-03-14T07:55:29Z",
            "user": {
                "id": 1,
                "username": "test",
                "email": "mail@mail.com",
                "created_at": "2025-03-12T13:00:21Z"
            }
        }]
*/
func (c *CommentStore) GetCommentOfUserOnPost(ctx context.Context, postId int64) ([]*Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id,users.email,users.created_at  FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
`

	rows, err := c.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	/*Something new in this section till end, its new with for next to get all the comments on that
	post id
	*/
	var comments []*Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.Username,
			&comment.User.ID,
			&comment.User.Email,
			&comment.User.Created,
		); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (c *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
INSERT INTO comments(post_id, user_id, content) values ($1,$2,$3)
returning id,created_at;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := c.db.QueryRowContext(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
