package users

import (
	"context"
	"errors"
	"go-crud-api/internal/config"
	"go-crud-api/pkg/password"
	"time"

	"go-crud-api/pkg/jwt"
)

// Service defines the user service.
type Service struct {
	repo   UserRepository
	config config.Config
}

// NewService creates a new user service.
func NewService(repo UserRepository, config config.Config) *Service {
	return &Service{repo: repo, config: config}
}

// Register creates a new user.
func (s *Service) Register(ctx context.Context, name, email, pass string) (*User, error) {
	hashedPassword, err := password.HashPassword(pass)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user", // Default role
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns access and refresh tokens.
func (s *Service) Login(ctx context.Context, email, pass string) (string, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err // Consider wrapping this error for better context
	}

	if !password.CheckPasswordHash(pass, user.PasswordHash) {
		return "", "", errors.New("invalid email or password")
	}

	accessTTL, _ := time.ParseDuration(s.config.AccessTokenTTL)
	refreshTTL, _ := time.ParseDuration(s.config.RefreshTokenTTL)

	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Role, s.config.JWTSecret, accessTTL, refreshTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// List returns all users.
func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}
