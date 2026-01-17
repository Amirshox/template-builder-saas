package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"template-builder-api/internal/queue"
	"template-builder-api/internal/repository"
	"template-builder-api/internal/service"
	"template-builder-api/pkg/db"
)

func main() {
	// 1. Init DB
	dbURL := "postgres://user:password@127.0.0.1:5433/template_builder?sslmode=disable"
	pool, err := db.NewPostgresDB(dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool)

	// 2. Init Asset Service (Internally inits MinIO, we don't need separate client here if we use the constructor)
	assetService, err := service.NewAssetService(repo, "localhost:9000", "minioadmin", "minioadmin")
	if err != nil {
		log.Fatal(err)
	}

	// 3. Init Renderer Service
	renderService := service.NewRenderService(repo, "http://localhost:3001")

	// 4. Init Queue
	q := queue.NewQueue("localhost:6380", "")

	log.Println("Worker started...")

	// 5. Start Consumer
	q.Consume(context.Background(), "workers-group", "worker-1", func(jobPayload queue.JobPayload) error {
		log.Printf("Processing Job: %s", jobPayload.JobID)

		ctx := context.Background()

		// 1. Update Status to Processing
		repo.UpdateJobStatus(ctx, jobPayload.JobID, "processing", nil, "")

		// 2. Call Renderer (Reusing Preview Logic but getting raw bytes)
		// Ideally RenderService should support generating generic IO Reader without version sometimes,
		// but here we use PreviewTemplate which fetches latest version by default if 0.
		pdfBytes, err := renderService.PreviewTemplate(ctx, jobPayload.TemplateID, 1)
		if err != nil {
			repo.UpdateJobStatus(ctx, jobPayload.JobID, "failed", nil, err.Error())
			return err
		}

		// 3. Upload to MinIO
		reader := strings.NewReader(string(pdfBytes)) // inefficient cast but OK for MVP
		filename := fmt.Sprintf("generated/%s.pdf", jobPayload.JobID)

		asset, err := assetService.UploadAsset(ctx, jobPayload.OrgID, reader, filename, int64(len(pdfBytes)), "application/pdf")
		if err != nil {
			repo.UpdateJobStatus(ctx, jobPayload.JobID, "failed", nil, "Failed to upload asset: "+err.Error())
			return err
		}

		// 4. Update Status to Completed
		repo.UpdateJobStatus(ctx, jobPayload.JobID, "completed", &asset.ID, "")

		log.Printf("Job Completed: %s", jobPayload.JobID)
		return nil
	})
}
