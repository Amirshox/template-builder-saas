package model

import (
	"time"

	"github.com/google/uuid"
)

type ValidationRules map[string]interface{}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Org struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Membership struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"` // owner, admin, editor, viewer
	CreatedAt time.Time `json:"created_at"`
}

type Template struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`   // layout, docx
	Status    string    `json:"status"` // active, archived
	CreatedAt time.Time `json:"created_at"`
}

type TemplateVersion struct {
	ID           uuid.UUID      `json:"id"`
	TemplateID   uuid.UUID      `json:"template_id"`
	Version      int            `json:"version"`
	Status       string         `json:"status"` // draft, published, archived
	TemplateJSON map[string]any `json:"template_json,omitempty"`
	SchemaJSON   map[string]any `json:"schema_json,omitempty"`
	DocxAssetID  *uuid.UUID     `json:"docx_asset_id,omitempty"`
	CreatedBy    *uuid.UUID     `json:"created_by,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	PublishedAt  *time.Time     `json:"published_at,omitempty"`
}
