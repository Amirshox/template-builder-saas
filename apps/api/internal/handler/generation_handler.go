package handler

import (
	"fmt"
	"net/http"
	"time"

	"template-builder-api/internal/model"
	"template-builder-api/internal/queue"
	"template-builder-api/internal/repository"
	"template-builder-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerationHandler struct {
	repo         repository.Repository
	queue        *queue.Queue
	assetService *service.AssetService
}

func NewGenerationHandler(repo repository.Repository, queue *queue.Queue, assetService *service.AssetService) *GenerationHandler {
	return &GenerationHandler{repo: repo, queue: queue, assetService: assetService}
}

func (h *GenerationHandler) GeneratePDF(c *gin.Context) {
	idStr := c.Param("id")
	templateID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	// Org ID from Auth
	orgID := c.MustGet("orgID").(uuid.UUID)

	jobID := uuid.New()
	job := &model.GenerationJob{
		ID:         jobID,
		OrgID:      orgID,
		TemplateID: templateID,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := h.repo.CreateJob(c.Request.Context(), job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job"})
		return
	}

	// Enqueue
	err = h.queue.EnqueueJob(c.Request.Context(), queue.JobPayload{
		JobID:      jobID,
		OrgID:      orgID,
		TemplateID: templateID,
		Data:       nil, // Can pass merge data here
	})
	if err != nil {
		h.repo.UpdateJobStatus(c.Request.Context(), jobID, "failed", nil, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue job"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"jobId": jobID})
}

func (h *GenerationHandler) GetJobStatus(c *gin.Context) {
	idStr := c.Param("id")
	jobID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	job, err := h.repo.GetJob(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	// If completed, fetch asset to get URL
	var response map[string]interface{}
	response = gin.H{
		"id":            job.ID,
		"orgId":         job.OrgID,
		"templateId":    job.TemplateID,
		"status":        job.Status,
		"outputAssetId": job.OutputAssetID,
		"errorMessage":  job.ErrorMessage,
		"createdAt":     job.CreatedAt,
		"updatedAt":     job.UpdatedAt,
	}

	if job.Status == "completed" && job.OutputAssetID != nil {
		url, err := h.assetService.GetDownloadURL(c.Request.Context(), *job.OutputAssetID)
		if err == nil {
			response["downloadUrl"] = url
		} else {
			fmt.Printf("GetDownloadURL Error: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, response)
}
