package users

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations.
 type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	List(ctx context.Context) ([]User, error)
	// TODO: Add Update, Delete methods as needed
}
