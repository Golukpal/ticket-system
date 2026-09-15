package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/Golukpal/ticket-system/internal/model"
	"github.com/Golukpal/ticket-system/internal/repository"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmail       = errors.New("email is required")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	users *repository.UserRepository
}

func NewAuthService(users *repository.UserRepository) *AuthService {
	return &AuthService{
		users: users,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
) (model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" {
		return model.User{}, ErrInvalidEmail
	}

	if len(password) < 8 {
		return model.User{}, ErrInvalidPassword
	}

	_, err := s.users.FindByEmail(ctx, email)

	if err == nil {
		return model.User{}, ErrEmailAlreadyExists
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		return model.User{}, fmt.Errorf("check existing user: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(
		ctx,
		email,
		string(passwordHash),
	)
	if err != nil {
		return model.User{}, fmt.Errorf("register user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return model.User{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.User{}, ErrInvalidCredentials
		}

		return model.User{}, fmt.Errorf("find user for login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return model.User{}, ErrInvalidCredentials
	}

	return user, nil
}