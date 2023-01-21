package service

import (
	"database/sql"
	"errors"
	"strings"

	"forum/internal/models"
	"forum/internal/repository"
)

type Post interface {
	CreatePost(post models.Post) error
	PostById(postID, UserID int) (models.Post, error)
	AllPosts(userID int) ([]models.Post, error)
	UsersPosts(userID int) ([]models.Post, error)
	PostsByCategory(userID int, category string) ([]models.Post, error)
	LikedPosts(userID int) ([]models.Post, error)
}

var (
	ErrEmptyPost = errors.New("can't create an empty post")
	ErrNoPost    = errors.New("post is not found")
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) CreatePost(post models.Post) error {
	if strings.TrimSpace(post.Content) == "" {
		return ErrEmptyPost
	}
	return s.repo.CreatePost(post)
}

func (s *PostService) AllPosts(userID int) ([]models.Post, error) {
	return s.repo.GetAllPosts(userID)
}

func (s *PostService) PostById(postID, UserID int) (models.Post, error) {
	posts, err := s.repo.GetPostById(postID, UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return posts, ErrNoPost
	} else if err != nil {
		return posts, err
	}

	return posts, nil
}

func (s *PostService) UsersPosts(userID int) ([]models.Post, error) {
	return s.repo.GetAllUserPosts(userID)
}

func (s *PostService) PostsByCategory(userID int, category string) ([]models.Post, error) {
	return s.repo.GetPostsByCategory(userID, category)
}

func (s *PostService) LikedPosts(userID int) ([]models.Post, error) {
	return s.repo.GetLikedPosts(userID)
}
