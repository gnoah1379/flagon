package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectGroup struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid;null"`
	OwnerID     uuid.UUID  `json:"owner_id" gorm:"type:uuid;not null"`
	Name        string     `json:"name" gorm:"not null"`
	Slug        string     `json:"slug" gorm:"unique;not null"`
	Description string     `json:"description" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ProjectGroupUser struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key"`
	GroupID   uuid.UUID `json:"group_id" gorm:"type:uuid;primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
