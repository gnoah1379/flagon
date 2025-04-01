package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user model
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Username  string    `json:"username" gorm:"not null"`
	Password  string    `json:"-" gorm:"not null"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	AvatarURL string    `json:"avatar_url" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSSO struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Provider   string    `json:"provider" gorm:"not null"`
	ProviderID string    `json:"provider_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
