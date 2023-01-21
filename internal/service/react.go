package service

import (
	"strconv"

	"forum/internal/models"
	"forum/internal/repository"
)

type Reaction interface {
	ReactToPost(postID, userID int, react string) error
	ReactToComment(commentID int, userID int, react string) (int, error)
}

type ReactionService struct {
	repo repository.Reaction
}

func NewReactionService(repo repository.Reaction) *ReactionService {
	return &ReactionService{
		repo: repo,
	}
}

func (s *ReactionService) ReactToPost(postID, userID int, react string) error {
	vote, err := strconv.Atoi(react)
	if err != nil {
		return err
	}

	reaction := models.Reaction{
		PostID: postID,
		UserID: userID,
		Vote:   vote,
	}
	return s.repo.CreateReactionPost(reaction)
}

func (s *ReactionService) ReactToComment(commentID int, userID int, react string) (int, error) {
	vote, err := strconv.Atoi(react)
	if err != nil {
		return 0, err
	}

	reaction := models.Reaction{
		CommentID: commentID,
		UserID:    userID,
		Vote:      vote,
	}
	return s.repo.CreateReactionComment(reaction)
}
