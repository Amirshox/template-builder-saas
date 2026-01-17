package model

import (
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Type        string    `json:"type"` // image, font, docx
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	S3Key       string    `json:"-"`
	URL         string    `json:"url,omitempty"` // Presigned URL for display
	CreatedAt   time.Time `json:"created_at"`
}
