package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectGroup struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parentId"`
	OwnerID     uuid.UUID  `json:"ownerId"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type ProjectGroupUser struct {
	UserID    uuid.UUID `json:"userId"`
	GroupID   uuid.UUID `json:"groupId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
