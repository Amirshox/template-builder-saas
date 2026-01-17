package handler

import (
	"net/http"
	"strconv"
	"template-builder-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreviewHandler struct {
	svc *service.RenderService
}

func NewPreviewHandler(svc *service.RenderService) *PreviewHandler {
	return &PreviewHandler{svc: svc}
}

func (h *PreviewHandler) PreviewTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	versionStr := c.Query("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		version = 1 // Default to 1 for MVP if not provided
	}

	pdfBytes, err := h.svc.PreviewTemplate(c.Request.Context(), id, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", `attachment; filename="preview.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
