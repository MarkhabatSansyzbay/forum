package service

import (
	"errors"
	"forum/internal/models"
	"forum/internal/repository"
	"strings"
)

type Commentary interface {
	CreateComment(comment models.Comment) error
	CommentsByPostID(ID int, userID int) ([]models.Comment, error)
}

var ErrEmptyComment = errors.New("can't create an empty comment")

type CommentService struct {
	repo repository.Commentary
}

func NewCommentService(repo repository.Commentary) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (s *CommentService) CreateComment(comment models.Comment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return ErrEmptyComment
	}
	return s.repo.CreateComment(comment)
}

func (s *CommentService) CommentsByPostID(ID int, userID int) ([]models.Comment, error) {
	return s.repo.CommentsByPostID(ID, userID)
}
