package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type Reaction interface {
	CreateReactionPost(reaction models.Reaction) error
	CreateReactionComment(reaction models.Reaction) (int, error)
}

type ReactionSqlite struct {
	db *sql.DB
}

func NewReactionSqlite(db *sql.DB) *ReactionSqlite {
	return &ReactionSqlite{
		db: db,
	}
}

func (s *ReactionSqlite) CreateReactionPost(reaction models.Reaction) error {
	queryInsert := `
        INSERT INTO REACTIONS (UserID, PostID, VOTE) VALUES ($1, $2, $3);
    `
	querySelect := `
		SELECT VOTE FROM REACTIONS WHERE UserID = $1 AND PostID = $2
	`
	var vote int
	if err := s.db.QueryRow(querySelect, reaction.UserID, reaction.PostID).Scan(&vote); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if _, err := s.db.Exec(queryInsert, reaction.UserID, reaction.PostID, reaction.Vote); err != nil {
			return err
		}
	} else {
		if _, err := s.db.Exec(`DELETE FROM REACTIONS WHERE UserID = $1 AND PostID = $2 AND VOTE = $3`, reaction.UserID, reaction.PostID, vote); err != nil {
			return err
		}
		if vote != reaction.Vote {
			if _, err := s.db.Exec(queryInsert, reaction.UserID, reaction.PostID, reaction.Vote); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ReactionSqlite) CreateReactionComment(reaction models.Reaction) (int, error) {
	var postID int
	if err := s.db.QueryRow(`SELECT PostID FROM COMMENTS WHERE ID = $1`, reaction.CommentID).Scan(&postID); err != nil {
		return postID, err
	}
	queryInsert := `
        INSERT INTO REACTIONS (UserID, CommentID, VOTE) VALUES ($1, $2, $3);
    `
	querySelect := `
		SELECT VOTE FROM REACTIONS WHERE UserID = $1 AND CommentID = $2
	`
	var vote int
	if err := s.db.QueryRow(querySelect, reaction.UserID, reaction.CommentID).Scan(&vote); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return postID, err
		}
		if _, err := s.db.Exec(queryInsert, reaction.UserID, reaction.CommentID, reaction.Vote); err != nil {
			return postID, err
		}
	} else {
		if _, err := s.db.Exec(`DELETE FROM REACTIONS WHERE UserID = $1 AND CommentID = $2 AND VOTE = $3`, reaction.UserID, reaction.CommentID, vote); err != nil {
			return postID, err
		}
		if vote != reaction.Vote {
			if _, err := s.db.Exec(queryInsert, reaction.UserID, reaction.CommentID, reaction.Vote); err != nil {
				return postID, err
			}
		}
	}

	return postID, nil
}
