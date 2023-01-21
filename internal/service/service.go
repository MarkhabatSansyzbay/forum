package service

import (
	"forum/internal/repository"
)

type Service struct {
	Authorization
	Post
	Commentary
	Reaction
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Post:          NewPostService(repo.Post),
		Commentary:    NewCommentService(repo.Commentary),
		Reaction:      NewReactionService(repo.Reaction),
	}
}
