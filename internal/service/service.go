package service

import (
	"forum/internal/repository"
)

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
