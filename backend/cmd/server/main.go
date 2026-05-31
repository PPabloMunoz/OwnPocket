package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/config"
	"github.com/ppablomunoz/ownpocket/backend/internal/handler"
)

var Version = "0.0.0-dev"

func main() {
	log.Printf("Starting OwnPocket %s", Version)
	cfg := config.LoadConfig()

	// Ensure the database directory exists
	dbDir := filepath.Dir(cfg.DBPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		log.Printf("Creating database directory: %s", dbDir)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
	}

	db, err := config.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database at %s: %v", cfg.DBPath, err)
	}

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	handler.SetupRoutes(api, db, cfg)

	// Register frontend routes
	handler.RegisterFrontend(r)

	log.Printf("Server running on :%d", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
