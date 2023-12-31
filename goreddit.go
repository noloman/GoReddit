package goreddit

import (
	"github.com/google/uuid"
)

type Thread struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}

type Post struct {
	ID            uuid.UUID `db:"id"`
	ThreadID      uuid.UUID `db:"thread_id"`
	Title         string    `db:"title"`
	Content       string    `db:"content"`
	Votes         int       `db:"votes"`
	CommentsCount int       `db:"comments_count"`
	ThreadTitle   string    `db:"thread_title"`
}

type Comment struct {
	ID      uuid.UUID `db:"id"`
	PostID  uuid.UUID `db:"post_id"`
	Content string    `db:"content"`
	Votes   int       `db:"votes"`
}

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}

// ThreadStore is an interface abstraction for interacting with threads
type ThreadStore interface {
	Thread(id uuid.UUID) (Thread, error)
	Threads() ([]Thread, error)
	CreateThread(t *Thread) error
	UpdateThread(t *Thread) error
	DeleteThread(id uuid.UUID) error
}

// PostStore is an interface abstraction for interacting with posts
type PostStore interface {
	Post(id uuid.UUID) (Post, error)
	PostsByThread(ThreadID uuid.UUID) ([]Post, error)
	Posts() ([]Post, error)
	CreatePost(t *Post) error
	UpdatePost(t *Post) error
	DeletePost(id uuid.UUID) error
}

// CommentStore is an interface abstraction for interacting with comments
type CommentStore interface {
	Comment(id uuid.UUID) (Comment, error)
	CommentsByPost(PostID uuid.UUID) ([]Comment, error)
	CreateComment(t *Comment) error
	UpdateComment(t *Comment) error
	DeleteComment(id uuid.UUID) error
}

type UserStore interface {
	User(id uuid.UUID) (User, error)
	UserByUsername(username string) (User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id uuid.UUID) error
}

// Store is a wrapper to help pass our DB stores to our app server using DI
type Store interface {
	ThreadStore
	PostStore
	CommentStore
	UserStore
}
