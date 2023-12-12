package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/noloman/goreddit"
)

type CommentStore struct {
	DB *sqlx.DB
}

func (c *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var cm goreddit.Comment
	if err := c.DB.Get(&cm, `SELECT * FROM comments where id = $1`, id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("Error retrieving comment: %w", err)
	}
	return cm, nil
}

func (c *CommentStore) CommentsByPost(postID uuid.UUID) ([]goreddit.Comment, error) {
	var cc []goreddit.Comment
	if err := c.DB.Select(&cc, `SELECT * FROM comments WHERE post_id = $1 ORDER BY votes DESC`, postID); err != nil {
		return nil, fmt.Errorf("Error retrieving comments in post: %w", err)
	}
	return cc, nil
}

func (c *CommentStore) CreateComment(cmt *goreddit.Comment) error {
	if err := c.DB.Get(cmt, `INSERT INTO comments VALUES ($1, $2, $3, $4) RETURNING *`,
		cmt.ID,
		cmt.PostID,
		cmt.Content,
		cmt.Votes); err != nil {
		return fmt.Errorf("Error creating comment: %w", err)
	}
	return nil
}

func (c *CommentStore) UpdateComment(t *goreddit.Comment) error {
	if err := c.DB.Get(t, `UPDATE comments SET post_id = $1, content = $2, votes = $3 WHERE id = $4 RETURNING *`,
		t.PostID,
		t.Content,
		t.Votes,
		t.ID); err != nil {
		return fmt.Errorf("Error updating comment: %w", err)
	}
	return nil
}

func (c *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := c.DB.Exec(`DELETE FROM comments where post_id = $1`, id); err != nil {
		return fmt.Errorf("Error deleting comment: %w", err)
	}
	return nil
}
