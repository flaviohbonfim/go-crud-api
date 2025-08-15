package users

import (
	"context"
	"errors"
	"go-crud-api/internal/config"
	"go-crud-api/pkg/password"
	"testing"
	

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	repo := new(MockUserRepository)
	cfg := config.Config{}
	service := NewService(repo, cfg)

	ctx := context.Background()
	name := "Test User"
	email := "test@example.com"
	pass := "password123"

	// Test case 1: Successful registration
	repo.On("Create", ctx, mock.AnythingOfType("*users.User")).Return(nil).Once()
	user, err := service.Register(ctx, name, email, pass)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "user", user.Role)
	assert.True(t, password.CheckPasswordHash(pass, user.PasswordHash))
	repo.AssertExpectations(t)

	// Test case 2: Repository returns an error
	repo.On("Create", ctx, mock.AnythingOfType("*users.User")).Return(errors.New("db error")).Once()
	user, err = service.Register(ctx, name, email, pass)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "db error")
	repo.AssertExpectations(t)
}

func TestUserService_Login(t *testing.T) {
	repo := new(MockUserRepository)
	cfg := config.Config{
		JWTSecret:       "testsecret",
		AccessTokenTTL:  "15m",
		RefreshTokenTTL: "7d",
	}
	service := NewService(repo, cfg)

	ctx := context.Background()
	email := "test@example.com"
	pass := "password123"

	hashedPassword, _ := password.HashPassword(pass)
	user := &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	// Test case 1: Successful login
	repo.On("FindByEmail", ctx, email).Return(user, nil).Once()
	accessToken, refreshToken, err := service.Login(ctx, email, pass)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	repo.AssertExpectations(t)

	// Test case 2: User not found
	repo.On("FindByEmail", ctx, email).Return(&User{}, gorm.ErrRecordNotFound).Once()
	accessToken, refreshToken, err = service.Login(ctx, email, pass)
	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	assert.Contains(t, err.Error(), gorm.ErrRecordNotFound.Error())
	repo.AssertExpectations(t)

	// Test case 3: Incorrect password
	wrongPassUser := &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: "wronghash", // Simulate wrong password
		Role:         "user",
	}
	repo.On("FindByEmail", ctx, email).Return(wrongPassUser, nil).Once()
	accessToken, refreshToken, err = service.Login(ctx, email, "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	repo.AssertExpectations(t)
}

func TestUserService_List(t *testing.T) {
	repo := new(MockUserRepository)
	cfg := config.Config{}
	service := NewService(repo, cfg)

	ctx := context.Background()

	// Test case 1: Successful list
	expectedUsers := []User{
		{ID: uuid.New(), Email: "user1@example.com"},
		{ID: uuid.New(), Email: "user2@example.com"},
	}
	repo.On("List", ctx).Return(expectedUsers, nil).Once()
	users, err := service.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	repo.AssertExpectations(t)

	// Test case 2: Repository returns an error
	repo.On("List", ctx).Return([]User{}, errors.New("db error")).Once()
	users, err = service.List(ctx)
	assert.Error(t, err)
	assert.Empty(t, users)
	assert.Contains(t, err.Error(), "db error")
	repo.AssertExpectations(t)
}
