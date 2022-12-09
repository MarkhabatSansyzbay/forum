package repository

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type AuthSqlite struct {
	db *sql.DB
}

func NewAuthSqlite(db *sql.DB) *AuthSqlite {
	return &AuthSqlite{
		db: db,
	}
}

func (s *AuthSqlite) CreateUser(user models.User) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (email, username, password) VALUES(?, ?, ?)", usersTable)
	result, err := s.db.Exec(query, user.Email, user.Username, user.Password)
	if err != nil {
		return 0, fmt.Errorf("create user -> error executing query: %s", err)
	}

	return result.LastInsertId()
}

func (s *AuthSqlite) GetUser(username string) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, username, password FROM %s WHERE username=?", usersTable)
	row := s.db.QueryRow(query, username)
	var user models.User
	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password); err != nil && err != sql.ErrNoRows {
		return user, err
	}

	return user, nil
}

func (s *AuthSqlite) CreateSession(session models.Session) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, token, expiration_date) VALUES(?, ?, ?)", sessionTable)
	res, err := s.db.Exec(query, session.UserId, session.Token, session.ExpirationDate)
	if err != nil {
		return 0, fmt.Errorf("create token -> error executing query: %s", err)
	}

	return res.LastInsertId()
}

func (s *AuthSqlite) GetSession(token string) (models.Session, error) {
	query := fmt.Sprintf("SELECT id, token, user_id, expiration_date FROM %s WHERE token=?", sessionTable)
	row := s.db.QueryRow(query, token)
	var session models.Session
	if err := row.Scan(&session.Id, &session.Token, &session.UserId, &session.ExpirationDate); err != nil && err != sql.ErrNoRows {
		return session, fmt.Errorf("get session -> error executing query: %s", err)
	}

	return session, nil
}

func (s *AuthSqlite) DeleteSession(token string) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE token=?", sessionTable)
	res, err := s.db.Exec(query, token)
	if err != nil {
		return 0, fmt.Errorf("delete session -> error executing query: %s", err)
	}

	return res.LastInsertId()
}

func (s *AuthSqlite) DeleteSessionByUserId(userId int64) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=?", sessionTable)
	res, err := s.db.Exec(query, userId)
	if err != nil {
		return 0, fmt.Errorf("delete session by user id -> error executing query: %s", err)
	}

	return res.LastInsertId()
}

func (s *AuthSqlite) UserByToken(token string) (models.User, error) {
	query := fmt.Sprintf("SELECT %s.Id, %s.Username, %s.Email FROM %s LEFT OUTER JOIN %s ON %s.id=%s.user_id AND token=?",
		usersTable, usersTable, usersTable, usersTable, sessionTable, usersTable, sessionTable)
	row := s.db.QueryRow(query, token)
	var user models.User
	if err := row.Scan(&user.Id, &user.Email, &user.Username); err != nil {
		return user, err
	}

	return user, nil
}
