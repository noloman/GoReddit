package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/noloman/goreddit"
)

type CommentStore struct {
	*sqlx.DB
}

func (c *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var cm goreddit.Comment
	if err := c.Get(cm, `SELECT * FROM comments where id = $1`, id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("Error retrieving comment: %w", err)
	}
	return cm, nil
}

func (c *CommentStore) CommentsByPost(PostID uuid.UUID) ([]goreddit.Comment, error) {
	var cm []goreddit.Comment
	if err := c.Get(cm, `SELECT * FROM comments WHERE post_id = $1`, PostID); err != nil {
		return nil, fmt.Errorf("Error retrieving comments in post: %w", err)
	}
	return cm, nil
}

func (c *CommentStore) CreateComment(t *goreddit.Comment) error {
	if err := c.Get(t, `INSERT INTO comment VALUES ($1, $2, $3) RETURNING *`,
		t.PostID,
		t.Content,
		t.Votes); err != nil {
		return fmt.Errorf("Error creating comment: %w", err)
	}
	return nil
}

func (c *CommentStore) UpdateComment(t *goreddit.Comment) error {
	if err := c.Get(t, `UPDATE comments SET post_id = $1, content = $2, votes = $3`,
		t.PostID, t.Content, t.Votes); err != nil {
		return fmt.Errorf("Error updating comment: %w", err)
	}
	return nil
}

func (c *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := c.Exec(`DELETE FROM comments where post_id = $1`, id); err != nil {
		return fmt.Errorf("Error deleting comment: %w", err)
	}
	return nil
}
