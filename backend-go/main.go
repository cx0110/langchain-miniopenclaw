package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"langchain-miniopenclaw-backend-go/config"
	"langchain-miniopenclaw-backend-go/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadSettings("") // empty string means look for config.yaml in current dir
	if err != nil {
		log.Printf("Warning: Failed to load config: %v. Using defaults and environment variables.", err)
		// Create a basic config with defaults
		cfg = &config.Settings{
			BackendDir:         ".",
			ProjectRoot:        "..",
			LLMProvider:        "zhipu",
			LLMModel:           "glm-5",
			ComponentCharLimit: 20000,
			TerminalTimeout:    30,
			ServerHost:         "0.0.0.0",
			ServerPort:         8002,
		}
	}

	// Override port from environment if specified
	if port := os.Getenv("PORT"); port != "" {
		if parsedPort, err := strconv.Atoi(port); err == nil {
			cfg.ServerPort = parsedPort
		}
	}

	// Set up Gin
	r := gin.New()

	// Add recovery middleware
	r.Use(gin.Recovery())

	// Add logger middleware
	r.Use(gin.Logger())

	// Set up handler with config
	h := handler.NewHandler(cfg)
	r.POST("/chat", h.ChatHandler)

	// Start server
	srvAddr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on %s", srvAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
