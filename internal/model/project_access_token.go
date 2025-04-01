package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessTokenType string

const (
	PersonalAccessToken AccessTokenType = "pat"
	GroupAccessToken    AccessTokenType = "gat"
	ProjectAccessToken  AccessTokenType = "pt"
)

type AccessToken struct {
	ID            uuid.UUID       `json:"id" gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Name          string          `json:"name" gorm:"not null"`
	Token         string          `json:"token" gorm:"unique;not null"`
	TokenType     AccessTokenType `json:"token_type" gorm:"not null"`
	GroupID       uuid.NullUUID   `json:"group_id" gorm:"type:uuid;null"`
	ProjectID     uuid.NullUUID   `json:"project_id" gorm:"type:uuid;null"`
	EnvironmentID uuid.NullUUID   `json:"environment_id" gorm:"type:uuid;null"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
