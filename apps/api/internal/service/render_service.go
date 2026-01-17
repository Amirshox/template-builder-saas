package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"template-builder-api/internal/repository"

	"github.com/google/uuid"
)

type RenderService struct {
	repo        repository.Repository
	rendererURL string
	client      *http.Client
}

func NewRenderService(repo repository.Repository, rendererURL string) *RenderService {
	return &RenderService{
		repo:        repo,
		rendererURL: rendererURL,
		client:      &http.Client{},
	}
}

type RenderRequest struct {
	TemplateJSON map[string]any `json:"templateJson"`
}

func (s *RenderService) PreviewTemplate(ctx context.Context, templateID uuid.UUID, version int) ([]byte, error) {
	// 1. Fetch Template Version
	// For MVP, if version is 0 (latest), we might need logic to find it.
	// Assuming handling explicit version for now.

	tmplVersion, err := s.repo.GetTemplateVersion(ctx, templateID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get template version: %w", err)
	}

	if tmplVersion.TemplateJSON == nil {
		return nil, fmt.Errorf("template has no layout json")
	}

	// 2. Prepare Request
	payload := RenderRequest{
		TemplateJSON: tmplVersion.TemplateJSON,
	}
	bodyBytes, _ := json.Marshal(payload)

	// 3. Call Renderer
	req, err := http.NewRequestWithContext(ctx, "POST", s.rendererURL+"/render", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call renderer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("renderer error: %s", string(body))
	}

	// 4. Return PDF bytes
	return io.ReadAll(resp.Body)
}
