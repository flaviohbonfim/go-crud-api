package products

import (
	"time"

	"github.com/google/uuid"
)

// Product represents the product model.
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"type:varchar(120);not null" json:"name" validate:"required,min=2,max=120"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:numeric(10,2);not null" json:"price" validate:"required,gte=0"`
	Stock       int       `gorm:"type:integer;not null" json:"stock" validate:"required,gte=0"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
