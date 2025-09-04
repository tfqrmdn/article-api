package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"article-api/internal/cache"
	"article-api/internal/config"
	"article-api/internal/database"
	"article-api/internal/handlers"
	"article-api/internal/migration"
	"article-api/internal/repository"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run database migrations
	if err := migration.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Redis cache (fallback to mock if Redis unavailable)
	var cacheService cache.CacheServiceInterface
	redisCache, err := cache.NewCacheService()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis (%v), using mock cache", err)
		cacheService = cache.NewMockCacheService()
	} else {
		cacheService = redisCache
		defer cacheService.Close()
	}

	// Initialize repository with cache
	articleRepo := repository.NewArticleRepository(db, cacheService)

	// Initialize handlers
	articleHandler := handlers.NewArticleHandler(articleRepo)

	// Setup routes
	router := http.NewServeMux()
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			articleHandler.ListArticles(w, r)
		case "POST":
			articleHandler.CreateArticle(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Use router directly without middleware
	handler := router

	// Create HTTP server with configurable settings
	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
