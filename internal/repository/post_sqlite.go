package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"

	"forum/internal/models"
)

type Post interface {
	CreatePost(post models.Post) error
	GetPostById(postID, UserID int) (models.Post, error)
	GetAllPosts(userID int) ([]models.Post, error)
	GetAllUserPosts(userID int) ([]models.Post, error)
	GetPostsByCategory(userID int, Category string) ([]models.Post, error)
	GetLikedPosts(userID int) ([]models.Post, error)
}

type PostSqlite struct {
	db *sql.DB
}

func NewPostSqlite(db *sql.DB) *PostSqlite {
	return &PostSqlite{
		db: db,
	}
}

const queryCountFeedback = `
	SELECT COUNT(*), (
		SELECT COUNT(*) FROM REACTIONS WHERE VOTE=-1 AND PostID = $1
	), (
		SELECT COUNT(*) FROM COMMENTS WHERE PostID = $1
	)
	FROM REACTIONS WHERE VOTE=1 AND PostID = $1
`

func (s *PostSqlite) CreatePost(post models.Post) error {
	query := `
        INSERT INTO POSTS (AuthorID, Title, Content) VALUES ($1, $2, $3)
    `

	res, err := s.db.Exec(query, post.AuthorID, post.Title, post.Content)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, category := range post.Categories {
		query := `
			INSERT INTO CATEGORIES (PostID, Category) VALUES ($1, $2)
		`
		if _, err := s.db.Exec(query, id, category); err != nil {
			return err
		}
	}

	for _, path := range post.ImagesPath {
		query := `
			INSERT INTO IMAGES (PostID, Image) VALUES ($1, $2)
		`
		if _, err := s.db.Exec(query, id, path); err != nil {
			return err
		}
	}

	return nil
}

func (s *PostSqlite) GetPostById(postID, UserID int) (models.Post, error) {
	query := `
		SELECT POSTS.ID, POSTS.AuthorID, POSTS.Title, POSTS.Content, USERS.Username 
		FROM POSTS INNER JOIN USERS ON USERS.ID=POSTS.AuthorID 
		WHERE POSTS.ID = $1
	`

	var post models.Post
	if err := s.db.QueryRow(query, postID).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Author); err != nil {
		return post, err
	}

	if err := s.db.QueryRow(queryCountFeedback, &post.ID).Scan(&post.LikeCount, &post.DislikeCount, &post.CommentCount); err != nil {
		return post, err
	}

	categories, err := s.getPostCategories(post.ID)
	if err != nil {
		return post, err
	}
	post.Categories = categories

	vote, err := s.getReactionToPost(UserID, post.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return post, err
		}
		vote = 0
	}
	post.Vote = vote

	images, err := s.getPostImages(post.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return post, err
		}
	}
	post.ImagesPath = images

	return post, nil
}

func (s *PostSqlite) GetAllPosts(userID int) ([]models.Post, error) {
	query := `
		SELECT POSTS.ID, POSTS.AuthorID, POSTS.Title, POSTS.Content, USERS.Username 
		FROM POSTS INNER JOIN USERS ON USERS.ID=POSTS.AuthorID
		ORDER BY POSTS.ID DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Author); err != nil {
			return posts, err
		}

		if err := s.db.QueryRow(queryCountFeedback, &post.ID).Scan(&post.LikeCount, &post.DislikeCount, &post.CommentCount); err != nil {
			return posts, err
		}

		categories, err := s.getPostCategories(post.ID)
		if err != nil {
			return posts, err
		}
		post.Categories = categories

		vote, err := s.getReactionToPost(userID, post.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return posts, err
			}
			vote = 0
		}
		post.Vote = vote

		images, err := s.getPostImages(post.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return posts, err
			}
		}
		post.ImagesPath = images

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostSqlite) GetAllUserPosts(userID int) ([]models.Post, error) {
	query := `
		SELECT POSTS.ID, POSTS.AuthorID, POSTS.Title, POSTS.Content, USERS.Username 
		FROM POSTS INNER JOIN USERS ON USERS.ID=POSTS.AuthorID AND USERS.ID=?
		ORDER BY POSTS.ID DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Author); err != nil {
			return posts, err
		}

		if err := s.db.QueryRow(queryCountFeedback, &post.ID).Scan(&post.LikeCount, &post.DislikeCount, &post.CommentCount); err != nil {
			return posts, err
		}

		categories, err := s.getPostCategories(post.ID)
		if err != nil {
			return posts, err
		}
		post.Categories = categories

		vote, err := s.getReactionToPost(userID, post.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return posts, err
			}
			vote = 0
		}
		post.Vote = vote

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostSqlite) GetPostsByCategory(UserID int, Category string) ([]models.Post, error) {
	query := `
		SELECT POSTS.ID, POSTS.AuthorID, POSTS.Title, POSTS.Content, USERS.Username 
		FROM POSTS INNER JOIN USERS ON USERS.ID=POSTS.AuthorID, CATEGORIES
		WHERE CATEGORIES.Category = $1 AND CATEGORIES.PostID=POSTS.ID
		ORDER BY POSTS.ID DESC
	`
	rows, err := s.db.Query(query, Category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Author); err != nil {
			return posts, err
		}
		if err := s.db.QueryRow(queryCountFeedback, &post.ID).Scan(&post.LikeCount, &post.DislikeCount, &post.CommentCount); err != nil {
			return posts, err
		}
		categories, err := s.getPostCategories(post.ID)
		if err != nil {
			return posts, err
		}
		post.Categories = categories
		vote, err := s.getReactionToPost(UserID, post.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return posts, err
			}
			vote = 0
		}
		post.Vote = vote
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostSqlite) GetLikedPosts(userID int) ([]models.Post, error) {
	query := `
		SELECT POSTS.ID, POSTS.AuthorID, POSTS.Title, POSTS.Content, USERS.Username
		FROM POSTS INNER JOIN USERS ON USERS.ID=POSTS.AuthorID, REACTIONS
		WHERE REACTIONS.PostID = POSTS.ID AND REACTIONS.VOTE = 1 AND REACTIONS.UserID = $1
		ORDER BY POSTS.ID DESC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		fmt.Println("*")
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Author); err != nil {
			return posts, err
		}
		if err := s.db.QueryRow(queryCountFeedback, &post.ID).Scan(&post.LikeCount, &post.DislikeCount, &post.CommentCount); err != nil {
			return posts, err
		}
		categories, err := s.getPostCategories(post.ID)
		if err != nil {
			return posts, err
		}
		post.Categories = categories

		post.Vote = 1

		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostSqlite) getPostCategories(postID int) ([]string, error) {
	const query = `
		SELECT Category FROM CATEGORIES WHERE PostID = $1
	`
	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return categories, err
		}

		categories = append(categories, category)
	}

	return categories, err
}

func (s *PostSqlite) getReactionToPost(userID int, postID int) (int, error) {
	query := `
		SELECT VOTE FROM REACTIONS WHERE UserID = $1 AND PostID = $2
	`

	var vote int
	if err := s.db.QueryRow(query, userID, postID).Scan(&vote); err != nil {
		return vote, err
	}

	return vote, nil
}

func (s *PostSqlite) getPostImages(postID int) ([]template.URL, error) {
	const query = `
		SELECT Image FROM IMAGES WHERE PostID = $1
	`
	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var images []template.URL
	for rows.Next() {
		var image template.URL
		if err := rows.Scan(&image); err != nil {
			return images, err
		}

		images = append(images, image)
	}

	return images, err
}
