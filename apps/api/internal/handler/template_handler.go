package handler

import (
	"net/http"
	"template-builder-api/internal/service"
	"template-builder-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TemplateHandler struct {
	svc *service.TemplateService
}

func NewTemplateHandler(svc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

type CreateTemplateRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=layout docx"`
	// OrgID string `json:"orgId" binding:"required"` // In real app, get from Context/Token
}

func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	// 1. Get Org ID from Context
	orgID := c.MustGet("orgID").(uuid.UUID)

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationError(err)})
		return
	}

	t, err := h.svc.CreateTemplate(c.Request.Context(), orgID, req.Name, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	orgID := c.MustGet("orgID").(uuid.UUID)

	templates, err := h.svc.ListTemplates(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

type CreateVersionRequest struct {
	TemplateJSON map[string]any `json:"templateJson"`
	SchemaJSON   map[string]any `json:"schemaJson"`
}

func (h *TemplateHandler) CreateVersion(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	// User ID from Auth
	userID := c.MustGet("userID").(uuid.UUID)

	var req CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	version, err := h.svc.CreateVersion(c.Request.Context(), templateID, userID, req.TemplateJSON, req.SchemaJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, version)
}

func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	t, err := h.svc.GetTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	// RBAC Check
	orgID := c.MustGet("orgID").(uuid.UUID)
	if t.OrgID != orgID {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *TemplateHandler) ListVersions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// TODO: Verify access to template (check OrgID) first.
	// For now assuming if they know ID they can list versions, or we fetch Template first.
	// Ideally Service handles this or we repeat the check.
	t, err := h.svc.GetTemplate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	orgID := c.MustGet("orgID").(uuid.UUID)
	if t.OrgID != orgID {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	versions, err := h.svc.ListVersions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}
