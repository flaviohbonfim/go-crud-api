package users

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user model.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:user_role;not null;default:user" json:"role"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
