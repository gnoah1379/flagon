package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Slug        string    `json:"slug" gorm:"unique;not null"`
	GroupID     uuid.UUID `json:"group_id" gorm:"type:uuid;not null"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectUser struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key"`
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
