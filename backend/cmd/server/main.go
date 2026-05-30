package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/config"
	"github.com/ppablomunoz/ownpocket/backend/internal/handler"
)

func main() {
	cfg := config.LoadConfig()

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

	log.Printf("Server running on :%d", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
