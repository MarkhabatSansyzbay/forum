package repository

import (
	"database/sql"
	"errors"
	"forum/internal/models"
)

type Commentary interface {
	CreateComment(comment models.Comment) error
	CommentsByPostID(ID int, userID int) ([]models.Comment, error)
}

type CommentSqlite struct {
	db *sql.DB
}

func NewCommentSqlite(db *sql.DB) *CommentSqlite {
	return &CommentSqlite{
		db: db,
	}
}

func (s *CommentSqlite) CreateComment(comment models.Comment) error {
	query := `
        INSERT INTO COMMENTS(AuthorID, PostID, Content) VALUES ($1, $2, $3)
    `

	if _, err := s.db.Exec(query, comment.UserID, comment.PostID, comment.Content); err != nil {
		return err
	}
	return nil
}

func (s *CommentSqlite) CommentsByPostID(ID int, userID int) ([]models.Comment, error) {
	query := `
		SELECT COMMENTS.ID, COMMENTS.AuthorID, COMMENTS.PostID, COMMENTS.Content, USERS.Username 
		FROM COMMENTS INNER JOIN USERS ON USERS.ID=COMMENTS.AuthorID 
		WHERE COMMENTS.PostID = $1
	`

	queryCount := `
		SELECT COUNT(*), (
			SELECT COUNT(*) FROM REACTIONS WHERE VOTE=-1 AND CommentID = $1
		)
		FROM REACTIONS WHERE VOTE=1 AND CommentID = $1
	`

	rows, err := s.db.Query(query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.Author); err != nil {
			return comments, err
		}

		if err := s.db.QueryRow(queryCount, &comment.ID).Scan(&comment.LikeCount, &comment.DislikeCount); err != nil {
			return comments, err
		}

		vote, err := s.getCommentReactions(userID, comment.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return comments, err
			}
			vote = 0
		}
		comment.Vote = vote

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return comments, err
	}

	return comments, nil
}

func (s *CommentSqlite) getCommentReactions(UserID int, CommentID int) (int, error) {
	query := `
		SELECT VOTE FROM REACTIONS WHERE UserID = $1 AND CommentID = $2
	`

	var vote int
	if err := s.db.QueryRow(query, UserID, CommentID).Scan(&vote); err != nil {
		return vote, err
	}

	return vote, nil
}
