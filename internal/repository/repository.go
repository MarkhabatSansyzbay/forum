package repository

import (
	"database/sql"
)

type Repository struct {
	Authorization
	Post
	Commentary
	Reaction
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSqlite(db),
		Post:          NewPostSqlite(db),
		Commentary:    NewCommentSqlite(db),
		Reaction:      NewReactionSqlite(db),
	}
}
