package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectFeatureFlag struct {
	ProjectID     uuid.UUID `json:"project_id" gorm:"type:uuid;primary_key"`
	FeatureID     uuid.UUID `json:"feature_id" gorm:"type:uuid;primary_key"`
	EnvironmentID uuid.UUID `json:"environment_id" gorm:"type:uuid;primary_key"`
	TargetGroupID uuid.UUID `json:"target_group_id" gorm:"type:uuid;primary_key"`
	Enabled       bool      `json:"enabled" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
