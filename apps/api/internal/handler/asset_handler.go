package handler

import (
	"net/http"
	"template-builder-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	svc *service.AssetService
}

func NewAssetHandler(svc *service.AssetService) *AssetHandler {
	return &AssetHandler{svc: svc}
}

func (h *AssetHandler) UploadAsset(c *gin.Context) {
	// Multipart form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// 1. Get Org ID from Auth
	orgID := c.MustGet("orgID").(uuid.UUID)

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	asset, err := h.svc.UploadAsset(c.Request.Context(), orgID, file, fileHeader.Filename, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}
