package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/config"
	"github.com/ppablomunoz/ownpocket/backend/internal/middleware"
	"github.com/ppablomunoz/ownpocket/backend/internal/service"
	"gorm.io/gorm"
)

// Handler holds dependencies for all routes
type Handler struct {
	db      *gorm.DB
	cfg     *config.Config
	service *service.Service
}

// NewHandler creates a new handler instance
func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:      db,
		cfg:     cfg,
		service: service.NewService(db),
	}
}

func SetupRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)

	api := r.Group("/v1")

	api.GET("/health", h.Health)
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret)))

	protected.GET("/test", h.Health)
}

func (h *Handler) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
