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

	// Accounts
	protected.GET("/accounts", h.GetAccounts)
	protected.POST("/accounts", h.CreateAccount)
	protected.GET("/accounts/:id", h.GetAccount)
	protected.PUT("/accounts/:id", h.UpdateAccount)
	protected.DELETE("/accounts/:id", h.DeleteAccount)

	// Transactions
	protected.GET("/transactions", h.GetTransactions)
	protected.POST("/transactions", h.CreateTransaction)
	protected.GET("/transactions/:id", h.GetTransaction)
	protected.PUT("/transactions/:id", h.UpdateTransaction)
	protected.DELETE("/transactions/:id", h.DeleteTransaction)

	// Categories
	protected.GET("/categories", h.GetCategories)
	protected.POST("/categories", h.CreateCategory)
	protected.PUT("/categories/:id", h.UpdateCategory)
	protected.DELETE("/categories/:id", h.DeleteCategory)

	// Budgets
	protected.GET("/budgets", h.GetBudgets)
	protected.POST("/budgets", h.CreateBudget)

	// Dashboard / Summary
	protected.GET("/dashboard/summary", h.GetDashboardSummary)
}

func (h *Handler) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
