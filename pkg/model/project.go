package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID     `json:"id"`
	Slug        string        `json:"slug"`
	GroupID     uuid.NullUUID `json:"groupId"`
	OwnerID     uuid.UUID     `json:"ownerId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type ProjectEnvironment struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Project     Project   `json:"project"`
}

type ProjectUser struct {
	UserID    uuid.UUID `json:"userId"`
	ProjectID uuid.UUID `json:"projectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
