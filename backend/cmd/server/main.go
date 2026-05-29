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

	fmt.Println("Successfully connected to data/app.db!")
	_ = db.Config

	r := gin.Default()

	// Middlewares
	r.Use(cors.Default())

	api := r.Group("/api")
	handler.SetupRoutes(api, db, cfg)

	log.Printf("Server running on :%d", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}
