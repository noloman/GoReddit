package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // side effect for Go Postgres driver for the database/sql package
)

// Store struct that contains the PostgreSQL stores
type Store struct {
	*ThreadStore
	*PostStore
	*CommentStore
}

// NewStore creates a concrete implementation of the Store struct with the dataSourceName
func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error connecting to database: %w", err)
	}
	return &Store{
		ThreadStore:  &ThreadStore{DB: db},
		PostStore:    &PostStore{DB: &sqlx.DB{}},
		CommentStore: &CommentStore{DB: db},
	}, nil
}
