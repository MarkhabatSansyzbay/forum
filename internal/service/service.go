package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int64, error)
	GetUser(user models.User) (models.User, error)
	SetSession(userId int64) (models.Session, error)
	CheckSession(token string) (int64, error)
	DeleteSession(token string) (int64, error)
	UserByToken(token string) (models.User, error)
}

type Post interface{}

type Service struct {
	Authorization
	Post
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
