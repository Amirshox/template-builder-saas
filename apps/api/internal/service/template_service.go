package service

import (
	"context"
	"fmt"
	"template-builder-api/internal/model"
	"template-builder-api/internal/repository"
	"time"

	"github.com/google/uuid"
)

type TemplateService struct {
	repo repository.Repository
}

func NewTemplateService(repo repository.Repository) *TemplateService {
	return &TemplateService{repo: repo}
}

func (s *TemplateService) CreateTemplate(ctx context.Context, orgID uuid.UUID, name string, tType string) (*model.Template, error) {
	if name == "" {
		return nil, fmt.Errorf("template name is required")
	}

	template := &model.Template{
		ID:        uuid.New(),
		OrgID:     orgID,
		Name:      name,
		Type:      tType,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

func (s *TemplateService) ListTemplates(ctx context.Context, orgID uuid.UUID) ([]model.Template, error) {
	return s.repo.ListTemplates(ctx, orgID)
}

func (s *TemplateService) GetTemplate(ctx context.Context, id uuid.UUID) (*model.Template, error) {
	return s.repo.GetTemplate(ctx, id)
}

func (s *TemplateService) ListVersions(ctx context.Context, templateID uuid.UUID) ([]model.TemplateVersion, error) {
	return s.repo.ListTemplateVersions(ctx, templateID)
}

func (s *TemplateService) CreateVersion(ctx context.Context, templateID uuid.UUID, userID uuid.UUID, templateJSON map[string]any, schemaJSON map[string]any) (*model.TemplateVersion, error) {
	maxVersion, err := s.repo.GetMaxVersion(ctx, templateID)
	if err != nil {
		return nil, err
	}
	newVersion := maxVersion + 1

	version := &model.TemplateVersion{
		ID:           uuid.New(),
		TemplateID:   templateID,
		Version:      newVersion,
		Status:       "draft",
		TemplateJSON: templateJSON,
		SchemaJSON:   schemaJSON,
		CreatedBy:    &userID,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateTemplateVersion(ctx, version); err != nil {
		return nil, err
	}
	return version, nil
}
