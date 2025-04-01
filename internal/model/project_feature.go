package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectFeature struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ProjectID    uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
	Name         string    `json:"name" gorm:"not null"`
	CategoryID   uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Description  string    `json:"description" gorm:"not null"`
	DefaultValue bool      `json:"default_value" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
