package repository

import (
	"database/sql"
)

type Post interface{}

type Repository struct {
	Authorization
	Post
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSqlite(db),
	}
}
