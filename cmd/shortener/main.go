package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drashti/url_shortner/internal/config"
	"github.com/drashti/url_shortner/internal/handlers"
	"github.com/drashti/url_shortner/internal/service"
	"github.com/drashti/url_shortner/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize PostgreSQL storage
	postgresStorage, err := storage.NewPostgresStorage(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
	}
	defer postgresStorage.Close()

	// Initialize Redis storage
	redisStorage, err := storage.NewRedisStorage(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	defer redisStorage.Close()

	// Initialize service
	shortenerService := service.NewShortenerService(postgresStorage, redisStorage)
	defer shortenerService.Close()

	// Initialize handler
	shortenerHandler := handlers.NewShortenerHandler(shortenerService)

	// Initialize Gin router
	router := gin.Default()

	// Register routes
	router.POST("/shorten", shortenerHandler.CreateShortURL)
	router.GET("/:shortCode", shortenerHandler.RedirectToURL)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
