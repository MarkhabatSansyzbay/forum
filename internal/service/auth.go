package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"forum/internal/models"
	"forum/internal/repository"
)

var (
	ErrInvalidData   = errors.New("invalid data to create a user")
	ErrNoUser        = errors.New("user doesn't exist")
	ErrWrongPassword = errors.New("wrong password")
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(newUser models.User) (int64, error) {
	// encrypt the password
	user, err := s.repo.GetUser(newUser.Username)
	// validation
	if err != nil {
		return 0, fmt.Errorf("create user -> get user error: %s", err)
	}
	if user.Username != "" {
		return 0, ErrInvalidData
	}
	return s.repo.CreateUser(newUser)
}

func (s *AuthService) GetUser(possibleUser models.User) (models.User, error) {
	user, err := s.repo.GetUser(possibleUser.Username)
	if err != nil {
		return user, err
	}

	if user.Username == "" {
		return user, ErrNoUser
	}

	if possibleUser.Password != user.Password {
		return user, ErrWrongPassword
	}
	return user, nil
}

func (s *AuthService) SetSession(userId int64) (models.Session, error) {
	s.repo.DeleteSessionByUserId(userId)

	token, err := s.generateToken()
	if err != nil {
		return models.Session{}, fmt.Errorf("set session -> error generating token: %s", err)
	}

	session := models.Session{
		UserId:         userId,
		Token:          token,
		ExpirationDate: time.Now().Add(time.Second * 45),
	}

	_, err = s.repo.CreateSession(session)
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (s *AuthService) DeleteSession(token string) (int64, error) {
	return s.repo.DeleteSession(token)
}

func (s *AuthService) CheckSession(token string) (int64, error) {
	session, err := s.repo.GetSession(token)
	if err != nil {
		return 0, err
	}

	if session.ExpirationDate.Before(time.Now()) {
		return s.repo.DeleteSession(token)
	}

	return 0, nil
}

func (s *AuthService) UserByToken(token string) (models.User, error) {
	return s.repo.UserByToken(token)
}

func (s *AuthService) generateToken() (string, error) {
	const tokenLength = 32
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// func (s *AuthService) encryptPassword(password string) string {

// }
