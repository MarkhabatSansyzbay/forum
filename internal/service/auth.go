package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"forum/internal/models"
	"forum/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Authorization interface {
	CreateUser(user models.User, isOauth2 bool) error
	SetSession(username, password string, isOauth2 bool) (models.Session, error)
	DeleteSession(token string) error
	UserByToken(token string) (models.User, error)
}

var (
	ErrNoUser        = errors.New("user doesn't exist")
	ErrWrongPassword = errors.New("wrong password")
	ErrUsernameTaken = errors.New("username is already taken")
	ErrEmailTaken    = errors.New("email address is already taken")
)

const sessionTime = time.Hour * 6

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(newUser models.User, isOauth2 bool) error {
	if _, err := s.repo.GetUser(newUser.Username); err != sql.ErrNoRows {
		if err == nil {
			return ErrUsernameTaken
		}
		return err
	}

	if !isOauth2 {
		if err := checkUserInfo(newUser); err != nil {
			return err
		}

		password, err := s.generatePasswordHash(newUser.Password)
		if err != nil {
			return err
		}

		newUser.Password = password
	}

	if err := s.repo.CreateUser(newUser); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrEmailTaken
		}
		return err
	}

	return nil
}

func (s *AuthService) SetSession(username, password string, isOauth2 bool) (models.Session, error) {
	user, err := s.checkUser(username, password, isOauth2)
	if err != nil {
		return models.Session{}, err
	}

	s.repo.DeleteSessionByUserId(user.ID)

	token, err := s.generateToken()
	if err != nil {
		return models.Session{}, fmt.Errorf("set session -> error generating token: %s", err)
	}

	session := models.Session{
		UserID:         user.ID,
		Token:          token,
		ExpirationDate: time.Now().Add(sessionTime),
	}

	if err = s.repo.CreateSession(session); err != nil {
		return session, fmt.Errorf("set session -> error creating session: %s", err)
	}

	return session, nil
}

func (s *AuthService) DeleteSession(token string) error {
	return s.repo.DeleteSession(token)
}

func (s *AuthService) UserByToken(token string) (models.User, error) {
	user, err := s.repo.UserByToken(token)
	if err != nil && err != sql.ErrNoRows {
		return user, nil
	}
	return user, nil
}

func (s *AuthService) checkUser(username, password string, isOauth2 bool) (models.User, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return user, ErrNoUser
	}

	if !isOauth2 {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return user, ErrWrongPassword
		}
	}

	return user, nil
}

func (s *AuthService) generateToken() (string, error) {
	const tokenLength = 32
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *AuthService) generatePasswordHash(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(pass), err
}
