package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectCategory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ProjectID   uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
