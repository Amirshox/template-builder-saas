package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"template-builder-api/internal/handler"
	"template-builder-api/internal/middleware"
	"template-builder-api/internal/queue"
	"template-builder-api/internal/repository"
	"template-builder-api/internal/service"
	"template-builder-api/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Init DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@127.0.0.1:5433/template_builder?sslmode=disable"
	}

	pool, err := db.NewPostgresDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	// 2. Init Layers
	repo := repository.NewPostgresRepository(pool)
	templateService := service.NewTemplateService(repo)

	// MinIO credentials (default from docker-compose)
	assetService, err := service.NewAssetService(repo, "localhost:9000", "minioadmin", "minioadmin")
	if err != nil {
		log.Fatalf("Failed to init asset service: %v", err)
	}
	// Ensure bucket exists on startup
	if err := assetService.EnsureBucket(context.Background()); err != nil {
		log.Printf("Warning: Failed to ensure bucket: %v", err)
	}

	renderService := service.NewRenderService(repo, "http://localhost:3001")
	authService := service.NewAuthService(repo)

	// Queue
	q := queue.NewQueue("localhost:6380", "")

	// 2.1 Init Handlers
	generationHandler := handler.NewGenerationHandler(repo, q, assetService)
	authHandler := handler.NewAuthHandler(authService)

	// 3. Init Router
	r := gin.Default()

	// Middleware (Simple CORS)
	// Middleware (Simple CORS)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Requested-With, Accept")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Auth Routes
	r.POST("/v1/register", authHandler.Register)
	r.POST("/v1/login", authHandler.Login)
	r.GET("/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Protected Routes
	api := r.Group("/v1")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// Template Handlers
		templateHandler := handler.NewTemplateHandler(templateService)
		api.POST("/templates", templateHandler.CreateTemplate)
		api.GET("/templates", templateHandler.ListTemplates)
		api.GET("/templates/:id", templateHandler.GetTemplate)
		api.GET("/templates/:id/versions", templateHandler.ListVersions)
		api.POST("/templates/:id/versions", templateHandler.CreateVersion)

		// Assets
		assetHandler := handler.NewAssetHandler(assetService)
		api.POST("/assets", assetHandler.UploadAsset)

		// Preview
		previewHandler := handler.NewPreviewHandler(renderService)
		api.POST("/templates/:id/preview", previewHandler.PreviewTemplate)

		// Generation
		api.POST("/templates/:id/generate", generationHandler.GeneratePDF)
		api.GET("/jobs/:id", generationHandler.GetJobStatus)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
