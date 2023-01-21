package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Authorization interface {
	CreateUser(user models.User) error
	GetUser(username, email string) (models.User, error)
	CreateSession(user models.Session) error
	GetSession(token string) (models.Session, error)
	DeleteSession(token string) error
	DeleteSessionByUserId(userID int) error
	UserByToken(token string) (models.User, error)
}

type AuthSqlite struct {
	db *sql.DB
}

func NewAuthSqlite(db *sql.DB) *AuthSqlite {
	return &AuthSqlite{
		db: db,
	}
}

func (s *AuthSqlite) CreateUser(user models.User) error {
	query := `
		INSERT INTO USERS (Username, Email, Password) VALUES ($1, $2, $3);
	`

	if _, err := s.db.Exec(query, user.Username, user.Email, user.Password); err != nil {
		return err
	}

	return nil
}

func (s *AuthSqlite) GetUser(username, email string) (models.User, error) {
	query := `
		SELECT ID, Username, Email, Password FROM USERS WHERE Username=$1 or Email = $2;
	`

	var user models.User

	if err := s.db.QueryRow(query, username, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		return user, err
	}

	return user, nil
}

func (s *AuthSqlite) CreateSession(session models.Session) error {
	query := `
		INSERT INTO SESSIONS (UserID, Token, ExpDate) VALUES ($1, $2, $3);
	`

	if _, err := s.db.Exec(query, session.UserID, session.Token, session.ExpirationDate); err != nil {
		return err
	}

	return nil
}

func (s *AuthSqlite) GetSession(token string) (models.Session, error) {
	query := `
		SELECT ID, UserID, Token, ExpDate FROM SESSIONS WHERE Token = ?;
	`
	var session models.Session
	if err := s.db.QueryRow(query, token).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpirationDate); err != nil {
		return session, err
	}
	return session, nil
}

func (s *AuthSqlite) DeleteSession(token string) error {
	query := `
		DELETE FROM SESSIONS WHERE Token = ?;
	`

	if _, err := s.db.Exec(query, token); err != nil {
		return err
	}
	return nil
}

func (s *AuthSqlite) DeleteSessionByUserId(userID int) error {
	query := `
		DELETE FROM SESSIONS WHERE UserID = ?;
	`

	if _, err := s.db.Exec(query, userID); err != nil {
		return err
	}
	return nil
}

func (s *AuthSqlite) UserByToken(token string) (models.User, error) {
	query := `
		SELECT USERS.ID, USERS.Username, USERS.Email, USERS.Password 
		FROM SESSIONS INNER JOIN USERS 
		ON USERS.ID = SESSIONS.UserID
		WHERE SESSIONS.Token = ?;
	`
	var user models.User
	if err := s.db.QueryRow(query, token).Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}
