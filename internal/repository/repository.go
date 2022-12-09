package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Authorization interface {
	CreateUser(user models.User) (int64, error)
	GetUser(username string) (models.User, error)
	CreateSession(user models.Session) (int64, error)
	GetSession(token string) (models.Session, error)
	DeleteSession(token string) (int64, error)
	DeleteSessionByUserId(userId int64) (int64, error)
	UserByToken(token string) (models.User, error)
}

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
