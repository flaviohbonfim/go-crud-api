package repository

import (
	"context"
	"go-crud-api/internal/domain/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository.
func NewGormUserRepository(db *gorm.DB) users.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(ctx context.Context, user *users.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormUserRepository) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	var user users.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*users.User, error) {
	var user users.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) List(ctx context.Context) ([]users.User, error) {
	var users []users.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}
