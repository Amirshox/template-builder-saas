package model

import (
	"time"

	"github.com/google/uuid"
)

type GenerationJob struct {
	ID            uuid.UUID  `json:"id"`
	OrgID         uuid.UUID  `json:"orgId"`
	TemplateID    uuid.UUID  `json:"templateId"`
	Status        string     `json:"status"` // pending, processing, completed, failed
	OutputAssetID *uuid.UUID `json:"outputAssetId,omitempty"`
	ErrorMessage  string     `json:"errorMessage,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
